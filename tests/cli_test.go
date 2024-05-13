package tests

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/smarty/assertions/should"
)

func TestCLI(t *testing.T) {
	if testing.Short() {
		t.Skip("long-running test")
	}
	cmd := exec.Command("go", "run", "github.com/mdwhatcott/calcy-apps/main/calc-cli", "-op", "+", "3", "4")

	out, err := cmd.CombinedOutput()

	should.So(t, err, should.BeNil)
	should.So(t, strings.TrimSpace(string(out)), should.Equal, "7")
}
