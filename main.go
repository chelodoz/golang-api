package main

import (
	handler "golang-api/handler"
	router "golang-api/router"
)

var (
	httpRouter router.Router      = router.NewMuxRouter()
	msgHandler handler.MsgHandler = handler.NewMsgHandler()
)

func main() {
	const port string = ":8080"
	httpRouter.GET("/", msgHandler.GetMsgs)
	httpRouter.SERVE(port)
}
