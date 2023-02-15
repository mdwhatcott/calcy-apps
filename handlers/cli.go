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
	if len(args) < 2 {
		return fmt.Errorf("%w (you provided %d)", notEnoughArgumentsError, len(args))
	}

	a, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("%w: %w", invalidArgumentError, err)
	}

	b, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("%w: %w", invalidArgumentError, err)
	}

	_, err = fmt.Fprint(this.output, this.calc.Calculate(a, b))
	if err != nil {
		return fmt.Errorf("%w: %w", writeError, err)
	}

	return nil
}

var (
	invalidArgumentError    = errors.New("invalid arg")
	notEnoughArgumentsError = errors.New("two arguments are required")
	writeError              = errors.New("write error")
)
