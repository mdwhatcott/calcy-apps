package httpserver

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

type logger interface {
	Printf(string, ...any)
}
type httpServer interface {
	Serve(listener net.Listener) error
	Shutdown(ctx context.Context) error
}

type Server struct {
	logger          logger
	softContext     context.Context
	softShutdown    context.CancelFunc
	hardContext     context.Context
	hardShutdown    context.CancelFunc
	shutdownTimeout time.Duration
	network         string
	address         string
	ready           func(bool)
	server          httpServer
}

func New(
	ctx context.Context, logger logger,
	shutdownTimeout time.Duration,
	network string, address string,
	ready func(bool), handler http.Handler,
) *Server {
	softContext, softShutdown := context.WithCancel(ctx)
	hardContext, hardShutdown := context.WithCancel(ctx)
	return &Server{
		logger:          logger,
		softContext:     softContext,
		softShutdown:    softShutdown,
		hardContext:     hardContext,
		hardShutdown:    hardShutdown,
		shutdownTimeout: shutdownTimeout,
		network:         network,
		address:         address,
		ready:           ready,
		server:          &http.Server{Handler: handler},
	}
}

func (this *Server) Listen() {
	var waiter sync.WaitGroup
	waiter.Add(2)
	defer waiter.Wait()

	go this.listen(waiter.Done)
	go this.watchShutdown(waiter.Done)
}

func (this *Server) listen(done func()) {
	defer done()
	listener, err := new(net.ListenConfig).Listen(this.softContext, this.network, this.address)
	this.listenReady(err == nil)
	if err != nil {
		this.logger.Printf("[WARN] Failed to listen %v\n", err)
		return
	}

	err = this.serve(listener)
	if err != nil {
		this.logger.Printf("[WARN] Failed to serve %v\n", err)
		return
	}
}
func (this *Server) listenReady(success bool) {
	if this.ready != nil {
		this.ready(success)
		this.ready = nil
	}
}
func (this *Server) serve(listener net.Listener) error {
	this.logger.Printf("[INFO] Serving http %s:%s\n", this.network, this.address)
	err := this.server.Serve(listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (this *Server) watchShutdown(done func()) {
	defer func() {
		done()
		this.hardShutdown()
		this.logger.Printf("[INFO] Shutdown complete\n")
	}()
	<-this.softContext.Done()
	ctx, cancel := context.WithTimeout(this.hardContext, this.shutdownTimeout)
	defer cancel()
	this.logger.Printf("[INFO] Shutting down HTTP server\n")
	_ = this.server.Shutdown(ctx)
}

func (this *Server) Close() error {
	this.softShutdown()
	return nil
}
