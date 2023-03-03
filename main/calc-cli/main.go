package main

import (
	"flag"
	"log"
	"os"

	"github.com/mdwhatcott/calcy-apps/handlers"
	"github.com/mdwhatcott/calcy-lib/calcy"
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
	"+": calcy.Addition{},
	"-": calcy.Subtraction{},
	"*": calcy.Multiplication{},
	"/": calcy.Division{},
	"?": calcy.Bogus{Offset: 42},
}
