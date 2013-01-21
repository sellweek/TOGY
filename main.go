package main

import (
	"TOGY/config"
	"TOGY/manager"
	"TOGY/updater"
	"flag"
	"fmt"
)

var configPath = flag.String("config", "config.json", "The path to the local config file.")
var coldStart = flag.Bool("coldStart", false, "Download active broadcast, current config and terminate.")

func main() {
	flag.Parse()

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
	err := manager.Run(*configPath)
	if err != nil {
		panic(err)
	}
}
