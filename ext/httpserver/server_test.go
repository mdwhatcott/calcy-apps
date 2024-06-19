package httpserver

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestServerFixture(t *testing.T) {
	gunit.Run(new(ServerFixture), t)
}

type ServerFixture struct {
	*gunit.Fixture
	lock *sync.Mutex
}

func (this *ServerFixture) Setup() {
	this.lock = new(sync.Mutex)
}

func (this *ServerFixture) Printf(format string, args ...any) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Fixture.Printf(format, args...)
}

func (this *ServerFixture) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusTeapot)
	_, _ = writer.Write([]byte("Hello, world!"))
}

func (this *ServerFixture) ready(success bool) {
	this.So(success, should.BeTrue)
}

func (this *ServerFixture) Test() {
	var waiter sync.WaitGroup
	waiter.Add(1)
	ctx := context.WithValue(context.Background(), "test", this.Name())
	server := New(ctx, this, time.Millisecond, "tcp", "localhost:8080", this.ready, this)
	go func() {
		defer waiter.Done()
		server.Listen()
	}()
	request, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	response, _ := http.DefaultClient.Do(request)
	this.So(response.StatusCode, should.Equal, http.StatusTeapot)
	_ = server.Close()
	waiter.Wait()

	request, _ = http.NewRequest("GET", "http://localhost:8080/", nil)
	response, err := http.DefaultClient.Do(request)
	this.So(err, should.NotBeNil)
	this.So(response, should.BeNil)
}
