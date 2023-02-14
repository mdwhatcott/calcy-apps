package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/mdwhatcott/calcy-lib/calcy"
)

type CSVHandler struct {
	reader *csv.Reader
	writer *csv.Writer
	logger *log.Logger
}

func NewCSVHandler(input io.Reader, output io.Writer, logger *log.Logger) *CSVHandler {
	return &CSVHandler{
		reader: csv.NewReader(input),
		writer: csv.NewWriter(output),
		logger: logger,
	}
}

func (this *CSVHandler) Handle() (err error) {
	defer func() {
		this.writer.Flush()
		if err == nil {
			err = this.writer.Error()
			if err != nil {
				err = fmt.Errorf("%w: %w", csvWriteError, err)
			}
		}
	}()

	for {
		record, err := this.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%w: %w", csvReadError, err)
		}

		a, err := strconv.Atoi(record[0])
		if err != nil {
			this.logger.Println("invalid operator:", err)
			continue
		}

		b, err := strconv.Atoi(record[2])
		if err != nil {
			this.logger.Println("invalid operator:", err)
			continue
		}

		calculator, ok := calculators[record[1]]
		if !ok {
			this.logger.Println("unsupported operand:", record[1])
			continue
		}

		result := calculator.Calculate(a, b)

		err = this.writer.Write(append(record, fmt.Sprint(result)))
		if err != nil {
			return fmt.Errorf("%w: %w", csvWriteError, err)
		}
	}
	return nil
}

var (
	csvReadError  = errors.New("csv read error")
	csvWriteError = errors.New("csv write error")
)

var calculators = map[string]calcy.Calculator{
	"+": calcy.Addition{},
	"-": calcy.Subtraction{},
	"*": calcy.Multiplication{},
	"/": calcy.Division{},
	"?": calcy.Bogus{Offset: 42},
}
