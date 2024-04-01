/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package bot

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/corpix/uarand"
	"github.com/go-errors/errors"
	"github.com/z3ntl3/VidmolySpoof/globals"
	"github.com/z3ntl3/VidmolySpoof/models"
	"go.mongodb.org/mongo-driver/bson"
	"h12.io/socks"

	"github.com/antchfx/htmlquery"
)

func ObtainManifest(body io.Reader) (string, error) {
	doc, err := htmlquery.Parse(body)
	if err != nil {
		return "", err
	}

	tree, err := htmlquery.Query(doc, "/html/body/script[7]")
	if err != nil {
		return "", err
	}

	if tree == nil {
		return "", errors.New("tree not parsed")
	}

	manifest := htmlquery.InnerText(tree)

	if !strings.Contains(manifest, "m3u8") {
		return "", errors.New("no manifest found")
	}
	manifest = strings.Split(strings.Split(manifest, "player.setup(")[1], ");")[0]
	manifest = strings.Split(strings.Split(manifest, "sources: [{file:\"")[1], "\"}")[0]

	return manifest, nil
}

type Client struct {
	*http.Client
}

type ManifestLink = string

func NewClient(timeout time.Duration) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: timeout,
			Jar:     http.DefaultClient.Jar,
		},
	}
}

/*
to unveil underlying m3u8 master manifestation
*/
func (c *Client) UnveilManifest(url string) (*ManifestLink, error) {
	workerPool := make(chan struct {
		Err  error
		Link ManifestLink
	}, 500)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	go func(pool chan<- struct {
		Err  error
		Link ManifestLink
	},
	) {
		var err error
		var link string

		defer func(err_ *error, link_ *string) {
			pool <- struct {
				Err  error
				Link string
			}{
				Err:  *err_,
				Link: *link_,
			}
		}(&err, &link)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return
		}

		build_headers(req)
		res, err := c.Client.Do(req)
		if err != nil {
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return
		}

		if res.StatusCode != 200 {
			err = errors.Errorf("Status code '%d' with body %s", res.StatusCode, body)
			return
		}

		link, err = ObtainManifest(strings.NewReader(string(body)))
	}(workerPool)

	for {
		select {
		case task := <-workerPool:
			if task.Err != nil {
				continue
			}
			return &task.Link, nil
		case <-ctx.Done():
			v := <-workerPool
			return nil, v.Err
		}
	}
}

/*
To obtain data rapidly about the vidmoly stream
*/
func StreamCore(molyLink string) (*models.StreamData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	task := make(chan struct {
		Err error
		Res models.StreamData
	})

	go func(worker chan<- struct {
		Err error
		Res models.StreamData
	},
	) {
		var err error
		var res models.StreamData

		defer func(res_ *models.StreamData, err_ *error) {
			worker <- struct {
				Err error
				Res models.StreamData
			}{
				Err: *err_,
				Res: *res_,
			}
		}(&res, &err)

		err = globals.MongoClient.Collection(models.StreamCol).FindOne(ctx, bson.M{
			"$match": bson.M{
				"vidmoly_alias": molyLink,
			},
		}).Decode(&res)
	}(task)

	select {
	case v := <-task:
		return &v.Res, v.Err
	case <-ctx.Done():
		{
			v := <-task
			return nil, v.Err
		}
	}
}

func build_headers(req *http.Request) {
	req.Header.Add("User-Agent", uarand.GetRandom())
	req.Header.Add("Cache-Control", "no-store")
}

func (c *Client) DelProxy() {
	c.Client.Transport = http.DefaultTransport
}

// socks 4/5 or http proxy
func (c *Client) SetProxy(proxyURI string) error {
	c.Client.Transport = &http.Transport{
		Dial: socks.Dial(proxyURI),
	}

	return nil
}
