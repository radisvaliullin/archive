package main

import (
	"github.com/radisvaliullin/test_task_15/pkg/loaner"
	"log"
)

func main() {

	// config and init loaner
	c := loaner.Config{
		InPath: "./task_desc/small",
		// InPath:  "./task_desc/large",
		OutPath: "./out",
	}
	l := loaner.New(c)

	// run loan handler
	err := l.Loan()
	if err != nil {
		log.Print("loan err: ", err)
		return
	}
}
