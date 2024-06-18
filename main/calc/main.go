package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mdwhatcott/calcy-apps/app/calculator"
	"github.com/mdwhatcott/calcy-apps/ext/httpstatus"
	HTTP "github.com/mdwhatcott/calcy-apps/http"
	"github.com/mdwhatcott/calcy-lib/calcy"
)

func main() {
	statusHandler := httpstatus.NewHandler(context.Background(), StaticOKHealthCheck{}, time.Second, time.Second, time.Second)
	go statusHandler.Listen()
	appHandler := calculator.NewHandler(
		calcy.Addition{},
		calcy.Subtraction{},
		calcy.Multiplication{},
		calcy.Division{},
	)
	endpoint := "localhost:8080"
	log.Println("Listening on", endpoint)
	err := http.ListenAndServe(endpoint, HTTP.Router(statusHandler, appHandler))
	if err != nil {
		log.Fatalln(err)
	}
}

type StaticOKHealthCheck struct{}

func (StaticOKHealthCheck) Status(ctx context.Context) error {
	// Usually this is where we would ping a database, or perform some operation to verify that the domain is in a functional state.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
