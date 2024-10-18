package tests

import (
	"io"
	"net/http"
	"net/http/httputil"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/smarty/assertions/should"
)

func TestHTTP(t *testing.T) {
	if testing.Short() {
		t.Skip("long-running test")
	}
	cmd := exec.Command("go", "run", "github.com/mdw-smarty/calc-apps/main/calc-http")
	err := cmd.Start()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	defer func() { _ = cmd.Process.Kill() }()

	time.Sleep(time.Millisecond * 500)

	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/add?a=3&b=4", nil)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	requestDump, err := httputil.DumpRequestOut(request, false)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	t.Log("Request:\n" + string(requestDump))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	t.Log("Response:\n" + string(responseDump))

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		t.Fatal("non-200 status:", response.Status)
	}

	out, err := io.ReadAll(response.Body)

	should.So(t, err, should.BeNil)
	should.So(t, strings.TrimSpace(string(out)), should.Equal, "7")
}
