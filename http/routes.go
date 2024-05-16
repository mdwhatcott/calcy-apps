package http

import (
	"net/http"

	"github.com/mdwhatcott/calcy-apps/app/contracts"
	"github.com/mdwhatcott/calcy-apps/http/inputs"

	"github.com/smarty/shuttle"
)

func Router(calculator contracts.Handler) http.Handler {
	h := http.NewServeMux()
	processor := func() shuttle.Processor { return NewProcessor(calculator) }
	h.Handle("/add",
		shuttle.NewHandler(
			shuttle.Options.InputModel(func() shuttle.InputModel { return inputs.NewAddition() }),
			shuttle.Options.Processor(processor),
		),
	)
	h.Handle("/sub",
		shuttle.NewHandler(
			shuttle.Options.InputModel(func() shuttle.InputModel { return inputs.NewSubtraction() }),
			shuttle.Options.Processor(processor),
		),
	)
	h.Handle("/mul",
		shuttle.NewHandler(
			shuttle.Options.InputModel(func() shuttle.InputModel { return inputs.NewMultiplication() }),
			shuttle.Options.Processor(processor),
		),
	)
	h.Handle("/div",
		shuttle.NewHandler(
			shuttle.Options.InputModel(func() shuttle.InputModel { return inputs.NewDivision() }),
			shuttle.Options.Processor(processor),
		),
	)
	return h
}
