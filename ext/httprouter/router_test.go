package httprouter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdwhatcott/calcy-apps/ext/httprouter"
	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestRouterFixture(t *testing.T) {
	gunit.Run(new(RouterFixture), t)
}

type RouterFixture struct {
	*gunit.Fixture
	router   http.Handler
	request  *http.Request
	response *httptest.ResponseRecorder
}

func (this *RouterFixture) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Fprintf(response, "%s:%s", request.Method, request.URL.Path)
}
func (this *RouterFixture) Setup() {
	var err error
	this.router, err = httprouter.New(
		httprouter.ParseRoute("GET", "/a/b/c/d", this),
	)
	this.So(err, should.BeNil)
	this.request = httptest.NewRequest("GET", "/", nil)
	this.response = httptest.NewRecorder()
}
func (this *RouterFixture) TestMatchingRoute_HTTP200() {
	this.request.URL.Path = "/a/b/c/d"

	this.router.ServeHTTP(this.response, this.request)

	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.Equal, "GET:/a/b/c/d")
}
func (this *RouterFixture) TestNoRoute_HTTP404() {
	this.request.URL.Path = "/nope"

	this.router.ServeHTTP(this.response, this.request)

	this.So(this.response.Code, should.Equal, http.StatusNotFound)
	this.So(this.response.Body.String(), should.Equal, "Not Found\n")
}
func (this *RouterFixture) TestWrongMethod_HTTP415() {
	this.request.URL.Path = "/a/b/c/d"
	this.request.Method = http.MethodPut

	this.router.ServeHTTP(this.response, this.request)

	this.So(this.response.Code, should.Equal, http.StatusMethodNotAllowed)
	this.So(this.response.Body.String(), should.Equal, "Method Not Allowed\n")
}
