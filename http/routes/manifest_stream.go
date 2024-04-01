/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package routes

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/grafov/m3u8"
	"github.com/z3ntl3/VidmolySpoof/bot"
)

/*
not done its test code will add better API
*/
func Manifest_Stream(ctx *gin.Context) {
	videoCtx := struct {
		URL string `validate:"required,vidmoly,max=300,min=5"`
	}{
		URL: ctx.Query("url"),
	}

	if err := binding.Validator.ValidateStruct(&videoCtx); err != nil {
		ctx.Abort()
		ctx.String(403, err.Error())
		return
	}

	client := bot.NewClient(time.Second * 5)
	manifest, err := client.UnveilManifest(videoCtx.URL)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithError(500, err)
		return
	}

	for k, v := range manifest.Headers {
		ctx.Header(k, strings.Join(v, ""))
	}

	master := m3u8.NewMasterPlaylist()
	{
		if err := master.Decode(*bytes.NewBufferString(manifest.Raw), true); err != nil {
			ctx.AbortWithError(500, err)
		}
	}

	for k, v := range master.Variants {
		uri, err := url.Parse(v.URI)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		uri.Host = "test.com"

		master.Variants[k].URI = uri.String()
	}

	if _, err := ctx.Writer.Write(master.Encode().Bytes()); err != nil {
		ctx.AbortWithError(500, err)
	}
}
