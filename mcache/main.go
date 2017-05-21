package main

import (
	"fmt"
	"test_task_11/mcache/mcache"
	"time"
)

//
func main() {

	s := mcache.NewStorage()
	s.Set("q", []string{"qwerty"}, time.Second*1)
	fmt.Printf("set - %+v\n", s)

	//s.Remove("q")
	//fmt.Printf("remove - %+v\n", s)

	s.Set("q", []string{"2qwerty"}, time.Second*8)
	fmt.Printf("set - %+v\n", s)

	time.Sleep(time.Second * 4)
	fmt.Printf("2s %+v\n", s)

	fmt.Println("keys", s.Keys())

	fmt.Println("mcache")
}
