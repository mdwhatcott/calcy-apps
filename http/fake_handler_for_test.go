package http

import (
	"context"

	"github.com/mdwhatcott/calcy-apps/app/commands"
)

type FakeApplicationHandler struct {
	ctx     func(context.Context)
	result  int
	err     error
	handled []any
}

func NewFakeApplicationHandler() *FakeApplicationHandler {
	return &FakeApplicationHandler{result: 42, ctx: func(context.Context) {}}
}

func (this *FakeApplicationHandler) Handle(ctx context.Context, messages ...any) {
	this.ctx(ctx)
	for _, message := range messages {
		switch command := message.(type) {
		case *commands.Add:
			this.handled = append(this.handled, *command)
			command.Result.C = this.result
			command.Result.Error = this.err
		case *commands.Subtract:
			this.handled = append(this.handled, *command)
			command.Result.C = this.result
			command.Result.Error = this.err
		case *commands.Multiply:
			this.handled = append(this.handled, *command)
			command.Result.C = this.result
			command.Result.Error = this.err
		case *commands.Divide:
			this.handled = append(this.handled, *command)
			command.Result.C = this.result
			command.Result.Error = this.err
		}
	}
}
