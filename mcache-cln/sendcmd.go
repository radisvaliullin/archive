package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"test_task_11/mcache/mcache"
)

//
func sendCommand(cmd *mcache.Command) ([]byte, error) {

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(cmd); err != nil {
		return nil, fmt.Errorf("command to jsom marshal err %v", err)
	}

	res, err := http.Post("http://0.0.0.0:7337/set", "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("send post err %v", err)
	}
	res.Close = true
	defer res.Body.Close()

	resJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body err %v", err)
	}

	return resJSON, nil
}
