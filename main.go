/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package main

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/z3ntl3/MolyRevProxy/db"
	httpServer "github.com/z3ntl3/MolyRevProxy/http"
	"github.com/z3ntl3/MolyRevProxy/http/routes"
	vld "github.com/z3ntl3/MolyRevProxy/http/validator"
)

func main() {
	db.ReadConfig()
	db.Connect(viper.GetStringMap("database")["uri"].(string))

	{
		runtime.GOMAXPROCS(viper.GetStringMap("threading")["cores"].(int))
	}

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
			HandleFuncs: []gin.HandlerFunc{routes.Manifest_Stream},
			Group:       server.Engine.Group("/stream"),
		},
	}

	server.RegisterRoutes(&routes)
	server.Run(":2000")
}
