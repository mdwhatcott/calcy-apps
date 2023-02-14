package main

import (
	"flag"
	"log"
	"os"

	"github.com/mdwhatcott/calcy-apps/handlers"
	"github.com/mdwhatcott/calcy-lib/calcy"
)

func main() {
	var op string
	flag.StringVar(&op, "op", "+", "Pick one: + - * / ?")
	flag.Parse()

	calculator, ok := calculators[op]
	if !ok {
		log.Fatalln("unsupported operand:", op)
	}
	handler := handlers.NewCLIHandler(calculator, os.Stdout)

	err := handler.Handle(flag.Args())
	if err != nil {
		log.Fatalln(err)
	}
}

var calculators = map[string]calcy.Calculator{
	"+": calcy.Addition{},
	"-": calcy.Subtraction{},
	"*": calcy.Multiplication{},
	"/": calcy.Division{},
	"?": calcy.Bogus{Offset: 42},
}
