package main

import (
	"fmt"
	"test_task_11/mcache/mcache"
	"time"
)

//
func main() {

	fmt.Println("start mcache sever")

	mcs := mcache.NewMCacheServer(":7337")
	mcs.Start()
	mcsErrs := mcs.GetSerErrChan()

	for {
		select {
		case err := <-mcsErrs:
			fmt.Println("mcache server err - ", err)
		case <-time.Tick(time.Second * 10):
			fmt.Println("mcache server heartbeat")

		}
	}
}
