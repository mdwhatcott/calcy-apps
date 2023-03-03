package handlers

import (
	"bytes"
	"errors"
	"testing"

	"github.com/mdwhatcott/calcy-lib/calcy"
)

func TestCLIHandler(t *testing.T) {
	var output bytes.Buffer
	handler := NewCLIHandler(calcy.Addition{}, &output)
	err := handler.Handle([]string{"1", "2"})
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if output.String() != "3" {
		t.Error("Want 3, got", output.String())
	}
}
func TestCLIHandler_unsupportedOperation(t *testing.T) {
	handler := NewCLIHandler(nil, nil)
	err := handler.Handle(nil)
	if !errors.Is(err, errUnsupportedOperation) {
		t.Error("unexpected error:", err)
	}
}
func TestCLIHandler_notEnoughArgumentsError(t *testing.T) {
	handler := NewCLIHandler(calcy.Addition{}, nil)
	err := handler.Handle([]string{""})
	if !errors.Is(err, errNotEnoughArguments) {
		t.Error("unexpected error:", err)
	}
}
func TestCLIHandler_invalidArgumentError(t *testing.T) {
	handler := NewCLIHandler(calcy.Addition{}, nil)
	err := handler.Handle([]string{"NaN", "2"})
	if !errors.Is(err, errInvalidArgument) {
		t.Error("unexpected error:", err)
	}
	err = handler.Handle([]string{"1", "NaN"})
	if !errors.Is(err, errInvalidArgument) {
		t.Error("unexpected error:", err)
	}
}
func TestCLIHandler_writeError(t *testing.T) {
	innerError := errors.New("write error")
	handler := NewCLIHandler(calcy.Addition{}, &ErringWriter{err: innerError})
	err := handler.Handle([]string{"1", "2"})
	if !errors.Is(err, errWrite) {
		t.Error("unexpected error:", err)
	}
	if !errors.Is(err, innerError) {
		t.Error("unexpected error:", err)
	}
}
