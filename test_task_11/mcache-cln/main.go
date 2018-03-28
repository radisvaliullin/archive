package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	srvAddr = "0.0.0.0:7337"
)

//
func main() {

	fmt.Println("mcache client")

	stdinReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("enter command: ")

		cmdStr, err := stdinReader.ReadString('\n')
		if err != nil {
			fmt.Println("command read err ", err)
			return
		}

		cmd, err := parseCommand(cmdStr)
		if err != nil {
			fmt.Println("command parse err ", err)
			continue
		}

		res, err := sendCommand(cmd)
		if err != nil {
			fmt.Println("command send to server err ", err)
			continue
		}
		if res.Success {
			if res.Result != nil {
				fmt.Println(*res.Result)
			} else {
				fmt.Println("")
			}
		} else {
			if res.Error != nil {
				fmt.Println(*res.Error)
			}
		}
	}
}
