package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mdwhatcott/calcy-apps/ext/should"
)

func TestHTTPHandler_404(t *testing.T) {
	response := httptest.NewRecorder()
	NewHTTPRouter().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/nope?a=1&b=2", nil))
	should.So(t, response.Code, should.Equal, http.StatusNotFound)
}
func TestHTTPHandler_200(t *testing.T) {
	testHTTP200(t, "/add?a=1&b=2", "3")
	testHTTP200(t, "/sub?a=5&b=3", "2")
	testHTTP200(t, "/mul?a=3&b=4", "12")
	testHTTP200(t, "/div?a=100&b=50", "2")
	testHTTP200(t, "/bog?a=1&b=2", "45")
}
func testHTTP200(t *testing.T, path, expectedResponseBody string) {
	t.Run(strings.TrimLeft(path, "/"), func(t *testing.T) {
		handler := NewHTTPRouter()
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		should.So(t, response.Code, should.Equal, http.StatusOK)
		actualResponseBody := strings.TrimSpace(response.Body.String())
		should.So(t, actualResponseBody, should.Equal, expectedResponseBody)
	})
}
func TestHTTPHandler_422_InvalidArgA(t *testing.T) {
	handler := NewHTTPRouter()
	request := httptest.NewRequest(http.MethodGet, "/add?a=NaN&b=2", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	should.So(t, response.Code, should.Equal, http.StatusUnprocessableEntity)
	should.So(t, strings.TrimSpace(response.Body.String()), should.Equal, "invalid 'a' parameter: [NaN]")
}
func TestHTTPHandler_422_InvalidArgB(t *testing.T) {
	handler := NewHTTPRouter()
	request := httptest.NewRequest(http.MethodGet, "/add?a=1&b=NaN", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	should.So(t, response.Code, should.Equal, http.StatusUnprocessableEntity)
	should.So(t, strings.TrimSpace(response.Body.String()), should.Equal, "invalid 'b' parameter: [NaN]")
}
