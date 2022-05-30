package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	const port string = ":8080"
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Up and running...")
	})
	log.Println("Server listening on port", port)
	log.Fatalln(http.ListenAndServe(port, router))
}
