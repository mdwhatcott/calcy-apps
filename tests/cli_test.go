package tests

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	if testing.Short() {
		t.Skip("long-running test")
	}
	cmd := exec.Command("go", "run", "github.com/mdwhatcott/calcy-apps/main/calc-cli", "-op", "+", "3", "4")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	output := strings.TrimSpace(string(out))
	if output != "7" {
		t.Error("Want 7, got", output)
	}
}
