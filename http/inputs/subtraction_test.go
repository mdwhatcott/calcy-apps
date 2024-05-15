package inputs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestSubtractionFixture(t *testing.T) {
	gunit.Run(new(SubtractionFixture), t)
}

type SubtractionFixture struct {
	*gunit.Fixture
	request *http.Request
	model   *Subtraction
}

func (this *SubtractionFixture) Setup() {
	this.request = httptest.NewRequest("GET", "/", nil)
	this.model = NewSubtraction()
}
func (this *SubtractionFixture) setQuery(key, value string) {
	query := this.request.URL.Query()
	query.Set(key, value)
	this.request.URL.RawQuery = query.Encode()
}

func (this *SubtractionFixture) TestBadParamA() {
	this.setQuery("a", "NaN")
	this.setQuery("b", "0")

	err := this.model.Bind(this.request)

	this.So(err, should.NotBeNil)
	this.So(this.model, should.Equal, NewSubtraction())
}
func (this *SubtractionFixture) TestBadParamB() {
	this.setQuery("a", "0")
	this.setQuery("b", "NaN")

	err := this.model.Bind(this.request)

	this.So(err, should.NotBeNil)
	this.So(this.model, should.Equal, NewSubtraction())
}
func (this *SubtractionFixture) TestHappy() {
	this.setQuery("a", "1")
	this.setQuery("b", "2")

	err := this.model.Bind(this.request)

	this.So(err, should.BeNil)
	this.So(this.model, should.Equal, &Subtraction{A: 1, B: 2})
}
