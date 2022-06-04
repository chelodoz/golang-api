package router

import (
	"golang-api/util"
	"net/http"
)

//Router is the interface to be implemented by the HTTP routers
type Router interface {
	GET(uri string, f func(http.ResponseWriter, *http.Request))
	POST(uri string, f func(http.ResponseWriter, *http.Request))
	PATCH(uri string, f func(http.ResponseWriter, *http.Request))
	DELETE(uri string, f func(http.ResponseWriter, *http.Request))
	SERVE(config util.Config)
}
