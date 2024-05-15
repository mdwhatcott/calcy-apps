package inputs

import "net/http"

type Division struct {
	A int
	B int
}

func NewDivision() *Division {
	return &Division{}
}
func (this *Division) Bind(request *http.Request) error {
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
