package inputs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestDivisionFixture(t *testing.T) {
	gunit.Run(new(DivisionFixture), t)
}

type DivisionFixture struct {
	*gunit.Fixture
	request *http.Request
	model   *Division
}

func (this *DivisionFixture) Setup() {
	this.request = httptest.NewRequest("GET", "/", nil)
	this.model = NewDivision()
}
func (this *DivisionFixture) setQuery(key, value string) {
	query := this.request.URL.Query()
	query.Set(key, value)
	this.request.URL.RawQuery = query.Encode()
}

func (this *DivisionFixture) TestBadParamA() {
	this.setQuery("a", "NaN")
	this.setQuery("b", "0")

	err := this.model.Bind(this.request)

	this.So(err, should.NotBeNil)
	this.So(this.model, should.Equal, NewDivision())
}
func (this *DivisionFixture) TestBadParamB() {
	this.setQuery("a", "0")
	this.setQuery("b", "NaN")

	err := this.model.Bind(this.request)

	this.So(err, should.NotBeNil)
	this.So(this.model, should.Equal, NewDivision())
}
func (this *DivisionFixture) TestHappy() {
	this.setQuery("a", "1")
	this.setQuery("b", "2")

	err := this.model.Bind(this.request)

	this.So(err, should.BeNil)
	this.So(this.model, should.Equal, &Division{A: 1, B: 2})
}
