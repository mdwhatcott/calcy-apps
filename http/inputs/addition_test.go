package inputs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestAdditionFixture(t *testing.T) {
	gunit.Run(new(AdditionFixture), t)
}

type AdditionFixture struct {
	*gunit.Fixture
	request *http.Request
	model   *Addition
}

func (this *AdditionFixture) Setup() {
	this.request = httptest.NewRequest("GET", "/", nil)
	this.model = NewAddition()
}
func (this *AdditionFixture) setQuery(key, value string) {
	query := this.request.URL.Query()
	query.Set(key, value)
	this.request.URL.RawQuery = query.Encode()
}

func (this *AdditionFixture) TestBadParamA() {
	this.setQuery("a", "NaN")
	this.setQuery("b", "0")

	err := this.model.Bind(this.request)

	this.So(err, should.NotBeNil)
	this.So(this.model, should.Equal, NewAddition())
}
func (this *AdditionFixture) TestBadParamB() {
	this.setQuery("a", "0")
	this.setQuery("b", "NaN")

	err := this.model.Bind(this.request)

	this.So(err, should.NotBeNil)
	this.So(this.model, should.Equal, NewAddition())
}
func (this *AdditionFixture) TestHappy() {
	this.setQuery("a", "1")
	this.setQuery("b", "2")

	err := this.model.Bind(this.request)

	this.So(err, should.BeNil)
	this.So(this.model, should.Equal, &Addition{A: 1, B: 2})
}
