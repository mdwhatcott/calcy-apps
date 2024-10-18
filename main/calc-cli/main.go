package main

import (
	"flag"
	"log"
	"os"

	"github.com/mdw-smarty/calc-apps/handlers"
	"github.com/mdw-smarty/calc-lib/calc"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	var op string
	flag.StringVar(&op, "op", "+", "Pick one: + - * / ?")
	flag.Parse()

	handler := handlers.NewCLIHandler(calculators[op], os.Stdout)

	err := handler.Handle(flag.Args())
	if err != nil {
		log.Fatalln(err)
	}
}

var calculators = map[string]handlers.Calculator{
	"+": calc.Addition{},
	"-": calc.Subtraction{},
	"*": calc.Multiplication{},
	"/": calc.Division{},
	"?": calc.Bogus{Offset: 42},
}
