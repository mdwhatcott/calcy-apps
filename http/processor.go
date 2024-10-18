package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mdw-smarty/calc-apps/app/commands"
	"github.com/mdw-smarty/calc-apps/app/contracts"
	"github.com/mdw-smarty/calc-apps/http/inputs"
	"github.com/mdw-smarty/calc-apps/http/views"

	"github.com/smarty/shuttle"
)

type Processor struct {
	handler contracts.Handler
}

func NewProcessor(handler contracts.Handler) *Processor {
	return &Processor{handler: handler}
}
func (this *Processor) Process(ctx context.Context, v any) any {
	switch input := v.(type) {

	case *inputs.Addition:
		return this.add(ctx, input)

	case *inputs.Subtraction:
		return this.sub(ctx, input)

	case *inputs.Multiplication:
		return this.mul(ctx, input)

	case *inputs.Division:
		return this.div(ctx, input)

	default:
		return internalServerError
	}
}
func (this *Processor) add(ctx context.Context, input *inputs.Addition) any {
	command := &commands.Add{A: input.A, B: input.B}
	this.handler.Handle(ctx, command)
	if command.Result.Error != nil {
		return additionFailure
	}
	return views.Addition{A: input.A, B: input.B, C: command.Result.C}
}
func (this *Processor) sub(ctx context.Context, input *inputs.Subtraction) any {
	command := &commands.Subtract{A: input.A, B: input.B}
	this.handler.Handle(ctx, command)
	if command.Result.Error != nil {
		return subtractionFailure
	}
	return views.Subtraction{A: input.A, B: input.B, C: command.Result.C}
}
func (this *Processor) mul(ctx context.Context, input *inputs.Multiplication) any {
	command := &commands.Multiply{A: input.A, B: input.B}
	this.handler.Handle(ctx, command)
	if command.Result.Error != nil {
		return multiplicationFailure
	}
	return views.Multiplication{A: input.A, B: input.B, C: command.Result.C}
}
func (this *Processor) div(ctx context.Context, input *inputs.Division) any {
	command := &commands.Divide{A: input.A, B: input.B}
	this.handler.Handle(ctx, command)
	if command.Result.Error != nil {
		return divisionFailure
	}
	return views.Division{A: input.A, B: input.B, C: command.Result.C}
}

var (
	additionFailure       = applicationError("calculation:addition-error", "added")
	subtractionFailure    = applicationError("calculation:subtraction-error", "subtracted")
	multiplicationFailure = applicationError("calculation:multiplication-error", "multiplied")
	divisionFailure       = applicationError("calculation:division-error", "division")
	internalServerError   = shuttle.SerializeResult{
		StatusCode: http.StatusInternalServerError,
		Content:    http.StatusText(http.StatusInternalServerError),
	}
)

func applicationError(name, verbPastParticiple string) shuttle.SerializeResult {
	return shuttle.SerializeResult{
		StatusCode: http.StatusInternalServerError,
		Content: shuttle.InputError{
			Fields:  []string{"query:a", "query:b"},
			Name:    name,
			Message: fmt.Sprintf("The operands could not be %s", verbPastParticiple),
		},
	}
}
