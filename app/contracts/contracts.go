package contracts

import "context"

type Handler interface {
	Handle(context.Context, ...any)
}
