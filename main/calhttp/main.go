package main

import (
	"log"
	"net/http"

	"github.com/mdwhatcott/calcy-apps/handlers"
)

func main() {
	endpoint := "localhost:8080"
	log.Println("Listening on", endpoint)
	err := http.ListenAndServe(endpoint, handlers.NewHTTPRouter())
	if err != nil {
		log.Fatalln(err)
	}
}
