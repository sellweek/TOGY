package main

import (
	"flag"
	"fmt"
	"github.com/sellweek/TOGY/config"
	"github.com/sellweek/TOGY/manager"
	"github.com/sellweek/TOGY/updater"
)

var configPath = flag.String("config", "config.json", "The path to the local config file.")
var coldStart = flag.Bool("coldStart", false, "Download active broadcast, current config and terminate.")

func main() {
	flag.Parse()

	//If the coldStart flag is set, we download
	//the central config, broadcast and terminate.
	if *coldStart {
		conf, err := config.ColdStart(*configPath)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		err = updater.ColdStart(conf)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		return
	}

	err := manager.RunBroadcast(*configPath)
	if err != nil {
		panic(err)
	}
}
