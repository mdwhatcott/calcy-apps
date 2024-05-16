package inputs

import (
	"net/http"

	"github.com/smarty/shuttle"
)

type Multiplication struct {
	shuttle.BaseInputModel
	A int
	B int
}

func NewMultiplication() *Multiplication {
	return &Multiplication{}
}
func (this *Multiplication) Bind(request *http.Request) error {
	query := request.URL.Query()
	a, err := parseInteger(query, "a")
	if err != nil {
		return err
	}
	this.A = a
	b, err := parseInteger(query, "b")
	if err != nil {
		return err
	}
	this.B = b
	return nil
}
