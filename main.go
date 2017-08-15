package main

import "redis-cli/genhandcli"
import "log"
import "flag"

//
func main() {

	var getErrMode = flag.Bool("getErrors", false, "Enable get errors mode")
	flag.Parse()

	conf := &genhandcli.CliConf{
		Addr:       "localhost:6379",
		GetErrMode: *getErrMode,
	}
	log.Printf("app conf - %+v", conf)

	cln := genhandcli.NewClient(conf)

	if err := cln.Start(); err != nil {
		log.Fatal("app fatal")
	}
}
