package shuttle

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestShuttleFixture(t *testing.T) {
	gunit.Run(new(ShuttleFixture), t)
}

type ShuttleFixture struct {
	*gunit.Fixture
	request       *http.Request
	response      *httptest.ResponseRecorder
	bindErr       error
	processResult SerializeResult
	handler       http.Handler
}

func (this *ShuttleFixture) Setup() {
	this.request = httptest.NewRequest("GET", "/", nil)
	this.request = this.request.WithContext(context.WithValue(context.Background(), "test", this.Name()))
	this.response = httptest.NewRecorder()
	this.handler = NewHandler(func() InputModel { return this }, func() Processor { return this })
}
func (this *ShuttleFixture) Bind(request *http.Request) error {
	this.So(request.Context().Value("test"), should.Equal, this.Name())
	return this.bindErr
}
func (this *ShuttleFixture) Process(ctx context.Context, v any) any {
	this.So(ctx.Value("test"), should.Equal, this.Name())
	this.So(v, should.Equal, this)
	return this.processResult
}

func (this *ShuttleFixture) TestBindError() {
	expected := InputError{Fields: []string{"field1", "field2"}, Name: this.Name(), Message: "Message"}
	this.bindErr = expected

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Code, should.Equal, http.StatusBadRequest)
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	var actual InputError
	_ = json.Unmarshal(this.response.Body.Bytes(), &actual)
	this.So(actual, should.Equal, expected)
}
func (this *ShuttleFixture) TestApplicationError() {
	this.processResult = SerializeResult{StatusCode: http.StatusTeapot, Content: 42}
	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Code, should.Equal, http.StatusTeapot)
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	var actual any
	_ = json.Unmarshal(this.response.Body.Bytes(), &actual)
	this.So(actual, should.Equal, 42)
}
func (this *ShuttleFixture) TestApplicationSuccess() {
	this.handler.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Header().Get("Content-Type"), should.BeBlank)
	this.So(this.response.Body.Len(), should.Equal, 0)
}
