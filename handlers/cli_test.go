package handlers

import (
	"bytes"
	"errors"
	"testing"

	"github.com/mdwhatcott/calcy-apps/ext/should"
	"github.com/mdwhatcott/calcy-lib/calcy"
)

func TestCLIHandler(t *testing.T) {
	var output bytes.Buffer
	handler := NewCLIHandler(calcy.Addition{}, &output)
	err := handler.Handle([]string{"1", "2"})
	should.So(t, err, should.BeNil)
	should.So(t, output.String(), should.Equal, "3")
}
func TestCLIHandler_unsupportedOperation(t *testing.T) {
	handler := NewCLIHandler(nil, nil)
	err := handler.Handle(nil)
	should.So(t, err, should.Wrap, errUnsupportedOperation)
}
func TestCLIHandler_notEnoughArgumentsError(t *testing.T) {
	handler := NewCLIHandler(calcy.Addition{}, nil)
	err := handler.Handle([]string{""})
	should.So(t, err, should.Wrap, errNotEnoughArguments)
}
func TestCLIHandler_invalidArgumentError(t *testing.T) {
	handler := NewCLIHandler(calcy.Addition{}, nil)
	err := handler.Handle([]string{"NaN", "2"})
	should.So(t, err, should.Wrap, errInvalidArgument)
	err = handler.Handle([]string{"1", "NaN"})
	should.So(t, err, should.Wrap, errInvalidArgument)
}
func TestCLIHandler_writeError(t *testing.T) {
	innerError := errors.New("write error")
	handler := NewCLIHandler(calcy.Addition{}, &ErringWriter{err: innerError})
	err := handler.Handle([]string{"1", "2"})
	should.So(t, err, should.Wrap, errWrite)
	should.So(t, err, should.Wrap, innerError)
}
