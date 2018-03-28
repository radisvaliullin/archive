package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"test_task_9/dump/dbdump"
)

//
func main() {

	// config
	dconf := &dbdump.DumperConf{DumpGoLimit: 2}

	// config from json
	confFile, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal("app conf open err ", err)
	}
	err = json.Unmarshal(confFile, &dconf.Instances)
	if err != nil {
		log.Fatal("app conf json unmarshal err ", err)
	}

	// dumping
	d := dbdump.NewDumper(dconf)
	d.ToDump()
}
