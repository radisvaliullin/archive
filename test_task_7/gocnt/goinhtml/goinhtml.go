// Package goinhtml - supports tools to get web page html and get count counts Go string in html.
package goinhtml

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

// URLRes - url scan result
type URLRes struct {
	ResMsg string
	Cnt    int
}

// URLScan - concurrence search "Go" in web page
func URLScan(url string, res chan URLRes, limit chan bool, wg *sync.WaitGroup) {
	defer func() {
		<-limit
		wg.Done()
	}()

	//
	html, err := getHTML(url)
	if err != nil {
		fmt.Println("URLScan: getHTML error", err)
	}

	// go counts in html
	gocnt := strings.Count(html, "Go")

	// send result
	resMsg := fmt.Sprintf("Count for %v: %v", url, gocnt)
	res <- URLRes{ResMsg: resMsg, Cnt: gocnt}
}

// getHTML - return html by url
func getHTML(url string) (html string, err error) {

	// http request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get err", err)
		return
	}
	defer resp.Body.Close()

	// read body bytes
	htmlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("html bytes read err", err)
		return
	}
	html = string(htmlBytes)

	return
}
