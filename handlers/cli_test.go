package handlers

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/mdw-smarty/calc-lib/calc"
	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestCLIHandler(t *testing.T) {
	gunit.Run(new(CLIHandlerFixture), t)
}

type CLIHandlerFixture struct {
	*gunit.Fixture
}

func (this *CLIHandlerFixture) handle(output io.Writer, args ...string) error {
	return NewCLIHandler(calc.Addition{}, output).Handle(args)
}
func (this *CLIHandlerFixture) TestCLIHandler() {
	var output bytes.Buffer
	err := this.handle(&output, "1", "2")
	this.So(err, should.BeNil)
	this.So(output.String(), should.Equal, "3")
}
func (this *CLIHandlerFixture) TestUnsupportedOperation() {
	this.So(NewCLIHandler(nil, nil).Handle(nil), should.Wrap, errUnsupportedOperation)
}
func (this *CLIHandlerFixture) TestNotEnoughArgumentsError() {
	this.So(this.handle(nil, ""), should.Wrap, errNotEnoughArguments)
}
func (this *CLIHandlerFixture) TestInvalidArgumentError() {
	this.So(this.handle(nil, "NaN", "2"), should.Wrap, errInvalidArgument)
	this.So(this.handle(nil, "1", "NaN"), should.Wrap, errInvalidArgument)
}
func (this *CLIHandlerFixture) TestWriteError() {
	innerError := errors.New("write error")
	err := this.handle(&ErringWriter{err: innerError}, "1", "2")
	this.So(err, should.Wrap, errWrite)
	this.So(err, should.Wrap, innerError)
}
