// Package goinhtml - supports tools to get web page html and get count counts Go string in html.
package goinhtml

import (
	"fmt"
	"golang.org/x/net/html/charset"
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
	html := getHTML(url)

	// go counts in html
	gocnt := strings.Count(html, "Go")

	// send result
	resMsg := fmt.Sprintf("Count for %v: %v", url, gocnt)
	res <- URLRes{ResMsg: resMsg, Cnt: gocnt}
}

// getHTML - return html by url
func getHTML(url string) (html string) {

	// http request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get err", err)
		return
	}
	defer resp.Body.Close()

	// convert to UTF-8
	htmlReader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println("resp body to utf-8 err", err)
		return
	}

	htmlBytes, err := ioutil.ReadAll(htmlReader)
	if err != nil {
		fmt.Println("html bytes read err", err)
		return
	}
	html = string(htmlBytes)

	return
}
