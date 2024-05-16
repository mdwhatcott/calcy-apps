package http

import (
	"net/http"

	"github.com/mdwhatcott/calcy-apps/app/contracts"
	"github.com/mdwhatcott/calcy-apps/ext/shuttle"
	"github.com/mdwhatcott/calcy-apps/http/inputs"
)

func Router(calculator contracts.Handler) http.Handler {
	h := http.NewServeMux()
	processor := func() shuttle.Processor { return NewProcessor(calculator) }
	h.Handle("/add", shuttle.NewHandler(func() shuttle.InputModel { return inputs.NewAddition() }, processor))
	h.Handle("/sub", shuttle.NewHandler(func() shuttle.InputModel { return inputs.NewSubtraction() }, processor))
	h.Handle("/mul", shuttle.NewHandler(func() shuttle.InputModel { return inputs.NewMultiplication() }, processor))
	h.Handle("/div", shuttle.NewHandler(func() shuttle.InputModel { return inputs.NewDivision() }, processor))
	return h
}
