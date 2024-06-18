package http

import (
	"context"
	"net/http"
	"time"

	"github.com/mdwhatcott/calcy-apps/app/contracts"
	"github.com/mdwhatcott/calcy-apps/ext/httpstatus"
	"github.com/mdwhatcott/calcy-apps/http/inputs"

	"github.com/smarty/httprouter"
	"github.com/smarty/shuttle"
)

func Router(calculator contracts.Handler) http.Handler {
	processor := func() shuttle.Processor { return NewProcessor(calculator) }
	router, err := httprouter.New(
		httprouter.Options.Routes(
			httprouter.ParseRoute("GET", "/status",
				httpstatus.NewHandler(context.Background(), StaticOKHealthCheck{}, time.Second, time.Second, time.Second),
			),
			httprouter.ParseRoute("GET", "/add",
				shuttle.NewHandler(
					shuttle.Options.InputModel(func() shuttle.InputModel { return inputs.NewAddition() }),
					shuttle.Options.Processor(processor),
				),
			),
			httprouter.ParseRoute("GET", "/sub",
				shuttle.NewHandler(
					shuttle.Options.InputModel(func() shuttle.InputModel { return inputs.NewSubtraction() }),
					shuttle.Options.Processor(processor),
				),
			),
			httprouter.ParseRoute("GET", "/mul",
				shuttle.NewHandler(
					shuttle.Options.InputModel(func() shuttle.InputModel { return inputs.NewMultiplication() }),
					shuttle.Options.Processor(processor),
				),
			),
			httprouter.ParseRoute("GET", "/div",
				shuttle.NewHandler(
					shuttle.Options.InputModel(func() shuttle.InputModel { return inputs.NewDivision() }),
					shuttle.Options.Processor(processor),
				),
			),
		),
	)
	if err != nil {
		panic(err)
	}
	return router
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
