package handler

import (
	"fmt"
	"net/http"
)

type msg struct{}

type MsgHandler interface {
	GetMsgs(rw http.ResponseWriter, r *http.Request)
}

func NewMsgHandler() MsgHandler {
	return &msg{}
}

func (m *msg) GetMsgs(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Up and running...")
}
