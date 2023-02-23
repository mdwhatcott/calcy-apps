package tests

import (
	"io"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestHTTP(t *testing.T) {
	if testing.Short() {
		t.Skip("long-running test")
	}
	cmd := exec.Command("go", "run", "github.com/mdwhatcott/calcy-apps/main/calc-http")
	err := cmd.Start()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	defer func() { _ = cmd.Process.Kill() }()

	time.Sleep(time.Millisecond * 500)

	response, err := http.Get("http://localhost:8080/add?a=3&b=4")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		t.Fatal("non-200 status:", response.Status)
	}

	out, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	output := strings.TrimSpace(string(out))
	if output != "7" {
		t.Error("Want 7, got", output)
	}
}
