package dto

import (
	"encoding/json"
	"net/http"
)

func WriteResponse[T any](rw http.ResponseWriter, code int, data T) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(code)
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		panic(err)
	}
}
