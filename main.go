package main

import (
	"TOGY/config"
	"TOGY/control"
	"TOGY/updater"
	"TOGY/util"
)

const configPath = "config.json"

var presentationOn = false

func main() {
	conf, err := config.Get(configPath)
	if err != nil {
		panic(err)
	}
	conf.Log.Println("Loaded config.")

	for {
		scrExit := startScreenMgr(conf)
		confChan := startUpdateMgr(conf)

		<-confChan

		newConf, err := config.Get(configPath)
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
					if !presentationOn {
						broadcast, err := updater.GetCurrentBroadcast(".")
						if err != nil {
							c.Log.Println("Could not get current broadcast: ", err)
							continue
						}

						err = control.StartPresentation(c.PowerPoint, broadcast)
						if err != nil {
							c.Log.Println("Could not start presentation: ", err)
							util.Sleep(20)
							continue
						}
						presentationOn = true
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
					if presentationOn {
						control.KillPresentation()
						presentationOn = false
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

	go func() {
		for {
			restart := updater.Update(c)
			if restart {
				configChan <- true
				return
			}
			util.Sleep(c.UpdateInterval)
		}
	}()

	return
}
