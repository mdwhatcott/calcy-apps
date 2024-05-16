package main

import (
	"log"
	"net/http"

	"github.com/mdwhatcott/calcy-apps/app/calculator"
	HTTP "github.com/mdwhatcott/calcy-apps/http"
	"github.com/mdwhatcott/calcy-lib/calcy"
)

func main() {
	appHandler := calculator.NewHandler(
		calcy.Addition{},
		calcy.Subtraction{},
		calcy.Multiplication{},
		calcy.Division{},
	)
	endpoint := "localhost:8080"
	log.Println("Listening on", endpoint)
	err := http.ListenAndServe(endpoint, HTTP.Router(appHandler))
	if err != nil {
		log.Fatalln(err)
	}
}
