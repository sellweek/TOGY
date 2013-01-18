package manager

import (
	"time"
)

func scheduleManager(mgr *Manager, t <-chan time.Time) {
	for {
		select {
		case _ = <-t:
			if mgr.config.Broadcast() {
				err := mgr.startBroadcast()
				if err != nil {
					mgr.config.Log.Println("Error when starting broadcast: ", err)
					continue
				}
				mgr.config.Log.Println("Broadcast on")
			} else {
				err := mgr.stopBroadcast()
				if err != nil {
					mgr.config.Log.Println("Error when stopping broadcast: ", err)
					continue
				}
				mgr.config.Log.Println("Broadcast off")
			}

		case _ = <-mgr.scheduleChan:
			mgr.config.Log.Println("Schedule manager terminating")
			break
		}
	}
}
