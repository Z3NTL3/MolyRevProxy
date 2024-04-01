/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package bot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/corpix/uarand"
	"github.com/go-errors/errors"
	"github.com/z3ntl3/VidmolySpoof/models"
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
				fmt.Println(task.Err.Error())
				continue
			}
			return &task.Link, nil
		case <-ctx.Done():
			v := <-workerPool
			return nil, v.Err
		}
	}
}

func (c *Client) Stream(ctx models.Playlist, path string) (io.ReadCloser, error) {
	// TODO
	return nil, nil
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
