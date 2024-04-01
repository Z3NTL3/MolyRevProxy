/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package http

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
}

type (
	HTTP_METHOD = string
	HTTP_PATH   = string
)

type Routes struct {
	Method      HTTP_METHOD
	Path        HTTP_PATH
	HandleFuncs []gin.HandlerFunc
	Group       *gin.RouterGroup
}

func (s *Server) RegisterRoutes(routes *[]Routes) {
	for _, v := range *routes {
		if v.Group != nil {
			v.Group.Handle(v.Method, v.Path, v.HandleFuncs...)
			continue
		}
		s.Handle(v.Method, v.Path, v.HandleFuncs...)
	}
}
