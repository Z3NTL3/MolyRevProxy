/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package routes

import (
	"fmt"
	"time"

	"github.com/go-errors/errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/z3ntl3/VidmolySpoof/bot"
)

/*
 */
func HLS_Stream(ctx *gin.Context) {
	videoCtx := struct {
		Video string `validate:"required,vidmoly,max=300,min=5"`
	}{
		Video: ctx.Query("url"),
	}

	if err := binding.Validator.ValidateStruct(&videoCtx); err != nil {
		ctx.Abort()
		ctx.String(403, err.Error())
		return
	}

	client := bot.NewClient(time.Second * 5)
	mURL, err := client.UnveilManifest(videoCtx.Video)
	if err != nil {
		fmt.Println(errors.New(err).ErrorStack())
		return
	}

	fmt.Println(mURL)
	// todo look into readme.md
}
