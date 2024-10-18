package tests

import (
	"bytes"
	"os/exec"
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

func TestCSV(t *testing.T) {
	if testing.Short() {
		t.Skip("long-running test")
	}
	stdOut := &bytes.Buffer{}
	cmd := exec.Command("go", "run", "github.com/mdw-smarty/calc-apps/main/calc-csv")
	cmd.Stdin = strings.NewReader(inputCSV)
	cmd.Stdout = stdOut

	err := cmd.Run()

	should.So(t, err, should.BeNil)
	should.So(t, stdOut.String(), should.Equal, expectedOutputCSV)
}
