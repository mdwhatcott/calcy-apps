package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestHTTPHandler(t *testing.T) {
	gunit.Run(new(HTTPHandlerFixture), t)
}

type HTTPHandlerFixture struct {
	*gunit.Fixture
}

func (this *HTTPHandlerFixture) assertResponse(path string, code int, body string) {
	response := httptest.NewRecorder()
	NewHTTPRouter().ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
	this.So(response.Code, should.Equal, code)
	this.So(strings.TrimSpace(response.Body.String()), should.Equal, body)
}
func (this *HTTPHandlerFixture) Test404() {
	this.assertResponse("/nope?a=1&b=2", http.StatusNotFound, "404 page not found")
}
func (this *HTTPHandlerFixture) Test200() {
	this.assertResponse("/add?a=1&b=2", http.StatusOK, "3")
	this.assertResponse("/sub?a=5&b=3", http.StatusOK, "2")
	this.assertResponse("/mul?a=3&b=4", http.StatusOK, "12")
	this.assertResponse("/div?a=100&b=50", http.StatusOK, "2")
	this.assertResponse("/bog?a=1&b=2", http.StatusOK, "45")
}
func (this *HTTPHandlerFixture) Test422() {
	this.assertResponse("/add?a=NaN&b=2", http.StatusUnprocessableEntity, "invalid 'a' parameter: [NaN]")
	this.assertResponse("/add?a=1&b=NaN", http.StatusUnprocessableEntity, "invalid 'b' parameter: [NaN]")
}
