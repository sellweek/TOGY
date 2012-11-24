package main

import (
	"TOGY/config"
	"TOGY/control"
	"TOGY/updater"
	"TOGY/util"
	"flag"
)

var configPath = flag.String("config", "./config.json", "The path to the local config file.")

var broadcastProcess control.Broadcast = nil

func main() {
	flag.Parse()
	conf, err := config.Get(*configPath)
	if err != nil {
		panic(err)
	}
	conf.Log.Println("Loaded config.")

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

	go func() {
		for {
			select {
			case <-exitChan:
				return
			default:
				if c.Broadcast() {
					if broadcastProcess == nil {
						broadcast, err := updater.GetCurrentBroadcast(".")
						if err != nil {
							c.Log.Println("Could not get current broadcast: ", err)
							continue
						}

						broadcastProcess = control.NewPowerPoint(c.PowerPoint, broadcast)
						err = broadcastProcess.Start()
						if err != nil {
							c.Log.Println("Could not start presentation: ", err)
							broadcastProcess = nil
							util.Sleep(20)
							continue
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
					if broadcastProcess != nil {
						broadcastProcess.Kill()
						broadcastProcess = nil
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
	//so that we don not pass nil because of starting the update
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
