package handlers

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/smarty/assertions/should"
)

var inputCSV = strings.Join([]string{
	"1,+,2",
	"2,-,1",
	"NaN,+,2",
	"1,+,NaN",
	"1,nop,2",
	"3,*,4",
	"20,/,10",
	"4,?,23",
}, "\n")

var expectedOutputCSV = strings.Join([]string{
	"1,+,2,3",
	"2,-,1,1",
	"3,*,4,12",
	"20,/,10,2",
	"4,?,23,69",
	"",
}, "\n")

func TestCSVHandler(t *testing.T) {
	var logOutput bytes.Buffer
	logger := log.New(&logOutput, "", 0)
	input := strings.NewReader(inputCSV)
	var output bytes.Buffer
	handler := NewCSVHandler(input, &output, logger)

	err := handler.Handle()

	should.So(t, err, should.BeNil)
	should.So(t, output.String(), should.Equal, expectedOutputCSV)
	if t.Failed() {
		t.Log("Log Output:\n" + logOutput.String())
	}
}
func TestCSVHandler_ReadError(t *testing.T) {
	var logOutput bytes.Buffer
	logger := log.New(&logOutput, "", 0)
	innerErr := errors.New("boink")
	input := &ErringReader{err: innerErr}
	var output bytes.Buffer
	handler := NewCSVHandler(input, &output, logger)

	err := handler.Handle()

	should.So(t, err, should.Wrap, innerErr)
	should.So(t, err, should.Wrap, csvReadError)
}
func TestCSVHandler_WriteError(t *testing.T) {
	var logOutput bytes.Buffer
	logger := log.New(&logOutput, "", 0)
	input := strings.NewReader(inputCSV)
	innerErr := errors.New("boink")
	output := &ErringWriter{err: innerErr}
	handler := NewCSVHandler(input, output, logger)

	err := handler.Handle()

	should.So(t, err, should.Wrap, innerErr)
	should.So(t, err, should.Wrap, csvWriteError)
}

type ErringReader struct{ err error }

func (this *ErringReader) Read([]byte) (int, error) {
	return 0, this.err
}

type ErringWriter struct{ err error }

func (this *ErringWriter) Write([]byte) (int, error) {
	return 0, this.err
}
