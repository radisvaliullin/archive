package main

import (
	"bufio"
	"fmt"
	"os"
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
		fmt.Println(string(res))
	}
}
