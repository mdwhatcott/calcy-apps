package should

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type testingT interface {
	Helper()
	Error(...any)
}

type assertion func(actual any, expected ...any) error

func So(t testingT, actual any, assert assertion, expected ...any) bool {
	err := assert(actual, expected...)
	if err != nil {
		t.Helper()
		t.Error(err)
	}
	return err == nil
}

type negated struct{}

var NOT negated

func Equal(actual any, EXPECTED ...any) error {
	expected := EXPECTED[0]
	if equalTimes(actual, expected) {
		return nil
	}
	if reflect.DeepEqual(expected, actual) {
		return nil
	}
	return fmt.Errorf("\n"+
		"Got:  (%s) %v\n"+
		"Want: (%s) %v\n",
		reflect.TypeOf(actual), actual,
		reflect.TypeOf(expected), expected,
	)
}
func (negated) Equal(actual any, expected ...any) error {
	if Equal(actual, expected...) != nil {
		return nil
	}
	return fmt.Errorf("\n"+
		"Got:      %s\n"+
		"Unwanted: %s\n",
		format(actual),
		format(expected[0]),
	)
}

func BeTrue(actual any, _ ...any) error          { return Equal(actual, true) }
func BeFalse(actual any, _ ...any) error         { return Equal(actual, false) }
func BeNil(actual any, _ ...any) error           { return Equal(actual, nil) }
func (negated) BeNil(actual any, _ ...any) error { return NOT.Equal(actual, nil) }
func Wrap(actual any, expected ...any) error {
	if errors.Is(actual.(error), expected[0].(error)) {
		return nil
	}
	return fmt.Errorf("expected %v to wrap %v (but it didn't)", actual, expected[0])
}

func equalTimes(a, b any) bool {
	return isTime(a) && isTime(b) && a.(time.Time).Equal(b.(time.Time))
}
func isTime(v any) bool {
	_, ok := v.(time.Time)
	return ok
}

func format(v any) string {
	if isTime(v) {
		return fmt.Sprintf("%v", v)
	}
	if v == nil {
		return fmt.Sprintf("%v", v)
	}
	if t := reflect.TypeOf(v); t.Kind() <= reflect.Float64 {
		return fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("%#v", v)
}
