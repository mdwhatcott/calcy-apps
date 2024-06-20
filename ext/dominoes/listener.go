package dominoes

import (
	"context"
	"io"
)

type Listener interface {
	Listen()
}
type ListenCloser interface {
	Listener
	io.Closer
}
type logger interface {
	Printf(string, ...interface{})
}

type linkedListener struct {
	current  Listener
	next     Listener
	ctx      context.Context
	shutdown context.CancelFunc
	managed  []io.Closer
	logger   logger
}

func New(listeners []Listener, resources []io.Closer) ListenCloser {
	if len(listeners) == 0 {
		return nil
	}
	listener := listeners[0]
	if listener == nil {
		panic("nil listener")
	}
	listeners = listeners[1:]
	var managed []io.Closer
	if len(listeners) == 1 {
		managed = resources
	}
	ctx, shutdown := context.WithCancel(context.Background())
	return &linkedListener{
		ctx:      ctx,
		shutdown: shutdown,
		current:  listener,
		next:     New(listeners, resources),
		managed:  managed,
	}
}

func (this *linkedListener) Listen() {
	if this.isLast() {
		this.listen()
	} else {
		go this.listen()
		this.next.Listen()
	}

}
func (this *linkedListener) listen() {
	this.current.Listen()
	<-this.ctx.Done()
	_close(this.next)
	if this.isLast() {
		for _, managed := range this.managed {
			_close(managed)
		}
	}
}

func (this *linkedListener) Close() error {
	defer this.shutdown()
	_close(this.current)
	return nil
}

func (this *linkedListener) isLast() bool {
	return this.next == nil
}

func _close(v any) {
	closer, ok := v.(io.Closer)
	if ok && closer != nil {
		_ = closer.Close()
	}
}
