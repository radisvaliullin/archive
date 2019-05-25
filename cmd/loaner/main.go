package main

import (
	"log"
	"github.com/radisvaliullin/test_task_15/pkg/loaner"
)

func main() {

	c := loaner.Config{
		InPath: "./task_desc/large",
		OutPath: "./out",
	}
	l := loaner.New(c)
	err := l.Loan()
	if err != nil {
		log.Print("loan err: ", err)
	}
	
}
