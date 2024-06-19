package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mdwhatcott/calcy-apps/app/calculator"
	"github.com/mdwhatcott/calcy-apps/ext/httpserver"
	HTTP "github.com/mdwhatcott/calcy-apps/http"
	"github.com/mdwhatcott/calcy-lib/calcy"
	"github.com/smarty/httpstatus"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Llongfile)
	statusHandler := httpstatus.New(
		httpstatus.Options.Context(context.Background()),
		httpstatus.Options.HealthCheck(StaticOKHealthCheck{}),
		httpstatus.Options.ResourceName("calcy-context"),
		httpstatus.Options.DisplayName("calcy"),
		httpstatus.Options.HealthCheckTimeout(time.Second),
		httpstatus.Options.HealthCheckFrequency(time.Second),
		httpstatus.Options.ShutdownDelay(time.Second),
	)
	go statusHandler.Listen()

	appHandler := calculator.NewHandler(
		calcy.Addition{},
		calcy.Subtraction{},
		calcy.Multiplication{},
		calcy.Division{},
	)
	router := HTTP.Router(statusHandler, appHandler)
	server := httpserver.New(
		context.Background(),
		logger,
		time.Second,
		"tcp",
		"localhost:8080",
		func(bool) {},
		router,
	)
	server.Listen()
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
