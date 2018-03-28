package goinhtml

import (
	"sync"
	"testing"
)

// TestGoScan -
func TestGoScan(t *testing.T) {

	urls := []string{
		"https://ru.wikipedia.org/wiki/Go",
		"https://en.wikipedia.org/wiki/Go_(programming_language)",
	}

	// setup
	limitChan := make(chan bool, 2)
	resChan := make(chan URLRes, 2)
	wg := &sync.WaitGroup{}

	//
	for _, url := range urls {
		limitChan <- true
		wg.Add(1)
		go URLScan(url, resChan, limitChan, wg)
	}

	// waiting
	wg.Wait()

	if len(resChan) != 2 {
		t.Error("Get result error")
	}
	//
	close(resChan)
	for res := range resChan {
		if res.Cnt == 0 {
			t.Error("Cant find Go")
		}
	}
}
