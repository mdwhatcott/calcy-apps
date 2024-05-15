package inputs

import "net/http"

type Subtraction struct {
	A int
	B int
}

func NewSubtraction() *Subtraction {
	return &Subtraction{}
}
func (this *Subtraction) Bind(request *http.Request) error {
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
