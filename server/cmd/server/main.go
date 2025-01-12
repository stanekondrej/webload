package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stanekondrej/webload/server/internal/app/server"
)

const PORT = 50000

var address string = fmt.Sprintf(":%d", PORT)

func main() {
	log.Println("Starting webload backend server")

	h := server.NewHandler()

	http.HandleFunc("/query", h.QueryFunc)     // query data
	http.HandleFunc("/provide", h.ProvideFunc) // provide data

	log.Printf("Listening on port %d", PORT)
	log.Fatal(http.ListenAndServe(address, nil))
}
