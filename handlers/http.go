package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mdwhatcott/calcy-lib/calcy"
)

func NewHTTPRouter() http.Handler {
	h := http.NewServeMux()
	h.Handle("/add", NewHTTPHandler(calcy.Addition{}))
	h.Handle("/sub", NewHTTPHandler(calcy.Subtraction{}))
	h.Handle("/mul", NewHTTPHandler(calcy.Multiplication{}))
	h.Handle("/div", NewHTTPHandler(calcy.Division{}))
	h.Handle("/bog", NewHTTPHandler(calcy.Bogus{Offset: 42}))
	return h
}

type HTTPHandler struct {
	calculator calcy.Calculator
}

func NewHTTPHandler(calculator calcy.Calculator) http.Handler {
	return &HTTPHandler{calculator: calculator}
}

func (this *HTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	rawA := query.Get("a")
	a, err := strconv.Atoi(rawA)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(writer, "invalid 'a' parameter: [%s]", rawA)
		return
	}

	rawB := query.Get("b")
	b, err := strconv.Atoi(rawB)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(writer, "invalid 'b' parameter: [%s]", rawB)
		return
	}

	writer.WriteHeader(200)
	_, err = fmt.Fprintln(writer, this.calculator.Calculate(a, b))
	if err != nil {
		log.Println("Failed to write response:", err)
	}
}
