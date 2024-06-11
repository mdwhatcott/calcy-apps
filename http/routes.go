package http

import (
	"net/http"

	"github.com/mdwhatcott/calcy-apps/app/contracts"
	"github.com/mdwhatcott/calcy-apps/ext/httprouter"
	"github.com/mdwhatcott/calcy-apps/http/inputs"

	"github.com/smarty/shuttle"
)

func Router(calculator contracts.Handler) http.Handler {
	processor := func() shuttle.Processor { return NewProcessor(calculator) }
	router, err := httprouter.New(
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
	)
	if err != nil {
		panic(err)
	}
	return router
}
