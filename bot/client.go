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
	"github.com/z3ntl3/MolyRevProxy/globals"
	"github.com/z3ntl3/MolyRevProxy/models"
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

type ManifestCtx struct {
	Headers http.Header
	Raw     string
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: timeout,
			Jar:     http.DefaultClient.Jar,
		},
	}
}

/*
to unveil underlying m3u8 manifest
*/
func (c *Client) GetManifest(url string, init bool) (*ManifestCtx, error) {
	workerPool := make(chan struct {
		Err error
		Ctx ManifestCtx
	}, 5)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	go func(pool chan struct {
		Err error
		Ctx ManifestCtx
	}, master bool,
	) {
		for i := 0; i < cap(pool); i++ {
			go func() {
				var err error
				var ctx ManifestCtx

				defer func(err_ *error, ctx_ *ManifestCtx) {
					pool <- struct {
						Err error
						Ctx ManifestCtx
					}{
						Err: *err_,
						Ctx: *ctx_,
					}
				}(&err, &ctx)

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

				if !master {
					ctx.Headers = req.Header
					ctx.Raw = string(body)

					return
				}

				link, err := ObtainManifest(strings.NewReader(string(body)))
				if err != nil {
					return
				}

				data, err := c.read_manifest(req, link)
				if err != nil {
					return
				}

				ctx.Headers = req.Header
				ctx.Raw = data
			}()
		}
	}(workerPool, init)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case task := <-workerPool:
			if task.Err != nil {
				continue
			}
			return &task.Ctx, nil
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

func (c *Client) read_manifest(req *http.Request, link string) (string, error) {
	var err error
	var result string

	req, err = http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return result, err
	}

	build_headers(req)
	res, err := c.Client.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	if res.StatusCode != 200 {
		err = errors.Errorf("Status code '%d' with body %s", res.StatusCode, body)
		return result, err
	}

	result = string(body)
	return result, err
}

func build_headers(req *http.Request) {
	req.Header.Add("User-Agent", uarand.GetRandom())
	req.Header.Add("Cache-Control", "no-store")
	req.Header.Add("Origin", "https://vidmoly.to")
	req.Header.Add("Referer", "https://vidmoly.to/")
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
