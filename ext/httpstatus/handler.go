package httpstatus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type HealthCheck interface {
	Status(ctx context.Context) error
}

type Handler struct {
	state         *atomic.Uint32
	hardContext   context.Context
	softContext   context.Context
	shutdown      context.CancelFunc
	healthCheck   HealthCheck
	timeout       time.Duration
	frequency     time.Duration
	shutdownDelay time.Duration
}

func NewHandler(ctx context.Context, check HealthCheck, timeout, frequency, shutdownDelay time.Duration) *Handler {
	softContext, shutdown := context.WithCancel(ctx)
	state := new(atomic.Uint32)
	state.Store(stateStarting)
	return &Handler{
		state:         state,
		hardContext:   ctx,
		softContext:   softContext,
		shutdown:      shutdown,
		healthCheck:   check,
		timeout:       timeout,
		frequency:     frequency,
		shutdownDelay: shutdownDelay,
	}
}

func (this *Handler) ServeHTTP(response http.ResponseWriter, _ *http.Request) {
	state := this.state.Load()
	switch state {
	case stateStarting, stateFailing, stateStopping:
		response.WriteHeader(http.StatusServiceUnavailable)
	case stateHealthy:
		response.WriteHeader(http.StatusOK)
	}
	_, _ = fmt.Fprint(response, statusText(state))
}
func (this *Handler) Listen() {
	for {
		err := this.healthCheck.Status(this.softContext)
		if err == nil {
			this.state.Store(stateHealthy)
		} else if errors.Is(err, context.Canceled) {
			this.state.Store(stateStopping)
			return
		} else {
			this.state.Store(stateFailing)
		}
		time.Sleep(this.frequency)
	}
}
func (this *Handler) Close() error {
	this.shutdown()
	return nil
}

func statusText(state uint32) string {
	switch state {
	case stateStarting:
		return "Starting"
	case stateHealthy:
		return "Healthy"
	case stateStopping:
		return "Stopping"
	default:
		return "Failing"
	}
}

const (
	stateStarting = iota
	stateHealthy
	stateFailing
	stateStopping
)
