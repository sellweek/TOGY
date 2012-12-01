package main

import (
	"TOGY/config"
	"TOGY/control"
	"TOGY/updater"
	"TOGY/util"
	"flag"
	"fmt"
)

var configPath = flag.String("config", "config.json", "The path to the local config file.")
var coldStart = flag.Bool("coldStart", false, "Download active broadcast, current config and terminate.")

var broadcastProcess control.Broadcast = nil

func main() {
	flag.Parse()
	conf, err := config.Get(*configPath)
	if err != nil {
		panic(err)
	}
	conf.Log.Println("Loaded config.")
	if *coldStart {
		err = updater.ColdStart(conf)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		return
	}

	for {
		scrExit := startScreenMgr(conf)
		confChan := startUpdateMgr(conf)

		<-confChan

		newConf, err := config.Get(*configPath)
		if err != nil {
			conf.Log.Println("Could not load new configuration file: ", err)
			continue
		} else {
			conf = newConf
		}
		conf.Log.Println("Loaded new config.")

		scrExit <- true
	}

}

func startScreenMgr(c config.Config) chan bool {
	exitChan := make(chan bool)

	if broadcastProcess == nil {
		broadcast, err := updater.GetCurrentBroadcast(".")
		if err != nil {
			c.Log.Println("Could not get current broadcast: ", err)
			panic("Could not get current broadcast")
		}

		broadcastProcess = control.NewPowerPoint(c.PowerPoint, broadcast)
	}

	go func() {
		for {
			select {
			case <-exitChan:
				return
			default:
				if c.Broadcast() {
					if !broadcastProcess.Status() {
						err := broadcastProcess.Start()
						if err != nil {
							c.Log.Println("Could not start presentation: ", err)
							broadcastProcess = nil
						}
					}
					err := control.TurnScreenOn()
					if err != nil {
						c.Log.Println("Could not turn screen on: ", err)
					}
					c.Log.Println("The screen is on")

				} else {
					err := control.TurnScreenOff()
					if err != nil {
						c.Log.Println("Could not turn screen off: ", err)
						continue
					}
					if broadcastProcess.Status() {
						broadcastProcess.Kill()
					}
					c.Log.Println("The screen is off")
				}
			}
			util.Sleep(10)
		}
	}()

	return exitChan
}

func startUpdateMgr(c config.Config) (configChan chan bool) {
	configChan = make(chan bool)
	//We have to wait until the current presentation is started
	//so that we do not pass nil because of starting the update
	//before the broadcast.
	util.Sleep(30)

	go func() {
		for {
			newBP, restart := updater.Update(c, broadcastProcess)
			if restart {
				configChan <- true
				return
			}
			if newBP != nil {
				broadcastProcess = newBP
			}
			util.Sleep(c.UpdateInterval)
		}
	}()

	return
}
