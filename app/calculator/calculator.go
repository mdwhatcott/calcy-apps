package calculator

import (
	"context"
	"fmt"

	"github.com/mdwhatcott/calcy-apps/app/commands"
)

type Calculator interface {
	Calculate(a, b int) int
}

type Handler struct {
	add Calculator
	sub Calculator
	mul Calculator
	div Calculator
}

func NewHandler(add, sub, mul, div Calculator) *Handler {
	return &Handler{
		add: add,
		sub: sub,
		mul: mul,
		div: div,
	}
}

func (this *Handler) Handle(_ context.Context, messages ...any) {
	for _, message := range messages {
		switch command := message.(type) {

		case *commands.Add:
			command.Result.C = this.add.Calculate(command.A, command.B)

		case *commands.Subtract:
			command.Result.C = this.sub.Calculate(command.A, command.B)

		case *commands.Multiply:
			command.Result.C = this.mul.Calculate(command.A, command.B)

		case *commands.Divide:
			command.Result.C = this.div.Calculate(command.A, command.B)

		default:
			panic(fmt.Sprintf("unsupported command: %T", command))
		}
	}
}
