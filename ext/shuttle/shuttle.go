package shuttle

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type (
	InputModel interface {
		Bind(request *http.Request) error
	}
	InputError struct {
		Fields  []string `json:"fields,omitempty"`
		Name    string   `json:"name,omitempty"`
		Message string   `json:"message,omitempty"`
	}
)

func (this InputError) Error() string {
	return fmt.Sprintf("%s %s", this.Name, this.Message)
}

type (
	Processor interface {
		Process(ctx context.Context, v any) any
	}
	SerializeResult struct {
		StatusCode int
		Content    any
	}
)

func NewHandler(input func() InputModel, processor func() Processor) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		model := input()
		if err := model.Bind(request); err != nil {
			respond(response, http.StatusBadRequest, err)
		} else {
			respond(response, http.StatusOK, processor().Process(request.Context(), model))
		}
	})
}
func respond(response http.ResponseWriter, defaultCode int, result any) {
	serializeResult, ok := result.(SerializeResult)
	if ok {
		result = serializeResult.Content
		if serializeResult.StatusCode > 0 {
			defaultCode = serializeResult.StatusCode
		}
	}
	if result != nil {
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	response.WriteHeader(defaultCode)
	if result == nil {
		return
	}
	encoder := json.NewEncoder(response)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(result)
	if err != nil {
		log.Println(err)
	}
}
