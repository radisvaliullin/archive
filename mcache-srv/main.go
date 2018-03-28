package main

import (
	"fmt"
	"test_task_11/mcache-srv/mcache"
	"time"
)

var (
	srvAddr = "0.0.0.0:7337"
)

//
func main() {

	fmt.Println("start mcache-srv sever")

	mcs := mcache.NewMCacheServer(srvAddr)
	mcs.Start()
	mcsErrs := mcs.GetSerErrChan()

	for {
		select {
		case err := <-mcsErrs:
			fmt.Println("mcache-srv server err - ", err)
		case <-time.Tick(time.Second * 15):
			fmt.Println("mcache-srv server heartbeat")

		}
	}
}
