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
		t.Error("Unexpected error:", err)
	}
	if output.String() != "3" {
		t.Error("Want 3, got", output.String())
	}
}
func TestCLIHandlerParseError(t *testing.T) {
	handler := NewCLIHandler(calcy.Addition{}, nil)
	err := handler.Handle([]string{"NaN", "2"})
	if !errors.Is(err, invalidArgumentError) {
		t.Error("Unexpected error:", err)
	}
	err = handler.Handle([]string{"1", "NaN"})
	if !errors.Is(err, invalidArgumentError) {
		t.Error("Unexpected error:", err)
	}
}
