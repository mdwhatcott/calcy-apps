package handlers

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

type CLIHandler struct {
	calc   Calculator
	output io.Writer
}

func NewCLIHandler(calculator Calculator, output io.Writer) *CLIHandler {
	return &CLIHandler{calc: calculator, output: output}
}

func (this *CLIHandler) Handle(args []string) error {
	if this.calc == nil {
		return fmt.Errorf("%w", errUnsupportedOperation)
	}

	if len(args) < 2 {
		return fmt.Errorf("%w (you provided %d)", errNotEnoughArguments, len(args))
	}

	a, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("%w: %w", errInvalidArgument, err)
	}

	b, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("%w: %w", errInvalidArgument, err)
	}

	_, err = fmt.Fprint(this.output, this.calc.Calculate(a, b))
	if err != nil {
		return fmt.Errorf("%w: %w", errWrite, err)
	}

	return nil
}

var (
	errUnsupportedOperation = errors.New("unsupported operation")
	errInvalidArgument      = errors.New("invalid arg")
	errNotEnoughArguments   = errors.New("two arguments are required")
	errWrite                = errors.New("write error")
)
