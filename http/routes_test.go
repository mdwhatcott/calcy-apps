package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestRoutesFixture(t *testing.T) {
	gunit.Run(new(RoutesFixture), t)
}

type RoutesFixture struct {
	*gunit.Fixture
}

func (this *RoutesFixture) route(method, URL string) {
	handler := NewFakeApplicationHandler()
	request := httptest.NewRequest(method, URL, nil)
	response := httptest.NewRecorder()
	router := Router(handler)
	router.ServeHTTP(response, request)
	this.So(response.Code, should.Equal, http.StatusOK)
}
func (this *RoutesFixture) TestRoutes() {
	this.route(http.MethodGet, "/add?a=1&b=2")
	this.route(http.MethodGet, "/sub?a=1&b=2")
	this.route(http.MethodGet, "/mul?a=1&b=2")
	this.route(http.MethodGet, "/div?a=1&b=2")
}
