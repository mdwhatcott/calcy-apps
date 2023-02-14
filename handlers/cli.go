package handlers

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/mdwhatcott/calcy-lib/calcy"
)

type CLIHandler struct {
	calc   calcy.Calculator
	output io.Writer
}

func NewCLIHandler(calculator calcy.Calculator, output io.Writer) *CLIHandler {
	return &CLIHandler{calc: calculator, output: output}
}

func (this *CLIHandler) Handle(args []string) error {
	a, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("%w: %w", invalidArgumentError, err)
	}

	b, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("%w: %w", invalidArgumentError, err)
	}

	_, err = fmt.Fprint(this.output, this.calc.Calculate(a, b))
	return err
}

var invalidArgumentError = errors.New("invalid arg")
