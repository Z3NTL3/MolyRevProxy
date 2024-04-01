/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	httpServer "github.com/z3ntl3/VidmolySpoof/http"
	"github.com/z3ntl3/VidmolySpoof/http/routes"
	vld "github.com/z3ntl3/VidmolySpoof/http/validator"
)

func main() {
	server := httpServer.Server{Engine: gin.New()}

	validator_ := &vld.Validator{
		Validate: validator.New(),
	}

	{
		validator_.RegisterValidation("vidmoly", func(fl validator.FieldLevel) bool {
			value := fl.Field().Interface().(string)

			if !strings.Contains(value, "https://vidmoly") {
				return false
			}
			return true
		})
	}

	binding.Validator = validator_

	routes := []httpServer.Routes{
		{
			Method:      http.MethodGet,
			Path:        "/manifest",
			HandleFuncs: []gin.HandlerFunc{routes.HLS_Stream},
			Group:       server.Engine.Group("/stream"),
		},
	}

	server.RegisterRoutes(&routes)
	server.Run(":2000")
}
