package inputs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestMultiplicationFixture(t *testing.T) {
	gunit.Run(new(MultiplicationFixture), t)
}

type MultiplicationFixture struct {
	*gunit.Fixture
	request *http.Request
	model   *Multiplication
}

func (this *MultiplicationFixture) Setup() {
	this.request = httptest.NewRequest("GET", "/", nil)
	this.model = NewMultiplication()
}
func (this *MultiplicationFixture) setQuery(key, value string) {
	query := this.request.URL.Query()
	query.Set(key, value)
	this.request.URL.RawQuery = query.Encode()
}

func (this *MultiplicationFixture) TestBadParamA() {
	this.setQuery("a", "NaN")
	this.setQuery("b", "0")

	err := this.model.Bind(this.request)

	this.So(err, should.NotBeNil)
	this.So(this.model, should.Equal, NewMultiplication())
}
func (this *MultiplicationFixture) TestBadParamB() {
	this.setQuery("a", "0")
	this.setQuery("b", "NaN")

	err := this.model.Bind(this.request)

	this.So(err, should.NotBeNil)
	this.So(this.model, should.Equal, NewMultiplication())
}
func (this *MultiplicationFixture) TestHappy() {
	this.setQuery("a", "1")
	this.setQuery("b", "2")

	err := this.model.Bind(this.request)

	this.So(err, should.BeNil)
	this.So(this.model, should.Equal, &Multiplication{A: 1, B: 2})
}
