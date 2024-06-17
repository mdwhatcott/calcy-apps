package httpstatus

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestHandlerFixture(t *testing.T) {
	gunit.Run(new(HandlerFixture), t)
}

type HandlerFixture struct {
	*gunit.Fixture
	ctx          context.Context
	statusErrors []error
	statusChecks *atomic.Int32
	response     *httptest.ResponseRecorder
	handler      *Handler
}

func (this *HandlerFixture) Status(ctx context.Context) (result error) {
	this.So(ctx.Value("test"), should.Equal, this.Name())

	select {
	case <-ctx.Done():
		result = ctx.Err()
	default:
		checks := int(this.statusChecks.Load())
		if len(this.statusErrors) > 0 && checks >= len(this.statusErrors) {
			result = this.statusErrors[len(this.statusErrors)-1]
		} else if checks < len(this.statusErrors) {
			result = this.statusErrors[checks]
		}
	}
	this.statusChecks.Add(1)
	return result
}

func (this *HandlerFixture) Setup() {
	this.statusChecks = new(atomic.Int32)
	this.ctx = context.WithValue(context.Background(), "test", this.Name())
	this.handler = NewHandler(this.ctx, this, time.Second, time.Millisecond, time.Millisecond)
	this.response = httptest.NewRecorder()
}

func (this *HandlerFixture) TestUponInitialization_StatusStarting() {
	this.handler.ServeHTTP(this.response, nil)
	this.So(this.response.Code, should.Equal, http.StatusServiceUnavailable)
	this.So(this.response.Body.String(), should.Equal, "Starting")
}
func (this *HandlerFixture) TestAfterSuccessfulHealthCheck_StatusHealthy() {
	go this.handler.Listen()
	time.Sleep(time.Millisecond * 10)
	this.handler.ServeHTTP(this.response, nil)
	this.So(this.statusChecks.Load(), should.BeGreaterThanOrEqualTo, 1)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.Equal, "Healthy")
}
func (this *HandlerFixture) TestAfterUnsuccessfulHealthCheck_StatusFailing() {
	this.statusErrors = append(this.statusErrors, errors.New("boink"))
	go this.handler.Listen()
	time.Sleep(time.Millisecond * 10)
	this.handler.ServeHTTP(this.response, nil)
	this.So(this.statusChecks.Load(), should.BeGreaterThanOrEqualTo, 1)
	this.So(this.response.Code, should.Equal, http.StatusServiceUnavailable)
	this.So(this.response.Body.String(), should.Equal, "Failing")
}
func (this *HandlerFixture) TestHealthyThanFailing() {
	this.statusErrors = append(this.statusErrors, nil, errors.New("boink"))
	go this.handler.Listen()
	time.Sleep(time.Millisecond * 10)
	this.handler.ServeHTTP(this.response, nil)
	this.So(this.statusChecks.Load(), should.BeGreaterThanOrEqualTo, 2)
	this.So(this.response.Code, should.Equal, http.StatusServiceUnavailable)
	this.So(this.response.Body.String(), should.Equal, "Failing")
}
func (this *HandlerFixture) TestAfterShutdownSignal_Stopping() {
	_ = this.handler.Close()
	this.handler.Listen()
	this.handler.ServeHTTP(this.response, nil)
	this.So(this.statusChecks.Load(), should.BeGreaterThanOrEqualTo, 1)
	this.So(this.response.Code, should.Equal, http.StatusServiceUnavailable)
	this.So(this.response.Body.String(), should.Equal, "Stopping")
}
