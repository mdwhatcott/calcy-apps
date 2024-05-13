package gunit

import (
	"reflect"
	"strings"
	"testing"
)

func Run(t *testing.T, fixture any) {
	fixtureType := reflect.TypeOf(fixture)
	for m := 0; m < fixtureType.NumMethod(); m++ {
		method := fixtureType.Method(m).Name
		if strings.HasPrefix(method, "Test") {
			t.Run(method, func(t *testing.T) {
				fixture := reflect.New(fixtureType.Elem())
				fixture.Elem().FieldByName("Fixture").Set(reflect.ValueOf(&Fixture{T: t}))
				if fixture.MethodByName("Setup").IsValid() {
					fixture.MethodByName("Setup").Call(nil)
				}
				fixture.MethodByName(method).Call(nil)
			})
		}
	}
}

type Fixture struct {
	*testing.T
}

func (this *Fixture) So(actual any, assert assertion, expected ...any) bool {
	err := assert(actual, expected...)
	if err != nil {
		this.Helper()
		this.Error(err)
	}
	return err == nil
}
func (this *Fixture) Write(p []byte) (int, error) {
	this.Log(strings.TrimSpace(string(p)))
	return len(p), nil
}

type assertion func(actual any, expected ...any) error
