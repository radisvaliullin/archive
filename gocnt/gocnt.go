// Package main - gocnt app implements simple count counts the Go string in webpage.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"test_task_7/gocnt/goinhtml"
)

// main - program entry point
func main() {

	// app setup
	// limit URLScan goroutines
	limit := 5
	// limit chan
	limitChan := make(chan bool, limit)
	// wait all run scan goroutines
	scanWG := &sync.WaitGroup{}

	// read from stdin (pipe)
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("stdin read err", err)
	}
	// stdin slice to string
	stdinString := strings.TrimSpace(string(stdin))

	// scanning urls slice
	urls := strings.Split(stdinString, "\n")

	// urls scan result chan
	resChan := make(chan goinhtml.URLRes, len(urls))

	// run Go scanning by url
	for _, url := range urls {
		// blocking if limit full
		limitChan <- true
		//
		scanWG.Add(1)
		go goinhtml.URLScan(url, resChan, limitChan, scanWG)
	}

	// waiting scaning goroutines
	scanWG.Wait()

	// Print result
	result := ""
	// all Go counts
	total := 0

	close(resChan)
	for res := range resChan {
		result += res.ResMsg + "\n"
		total += res.Cnt
	}
	result += fmt.Sprintf("Total: %v", total)

	//
	fmt.Println(result)

}
