package manager

import (
	"time"
)

//scheduleManager receives signals from a Ticker chan,
//sending a message to turn the broadcast on or off,
//depending on times set in Config.
func scheduleManager(mgr *Manager, t <-chan time.Time) {
	for {
		mgr.config.Debug("Schedule manager iteration")
		select {
		case _ = <-t:
			if mgr.config.Broadcast() {
				mgr.config.Debug("Sending message to start the broadcast")
				err := mgr.startBroadcast()
				if err != nil {
					mgr.config.Error("Error when starting broadcast: %v", err)
					continue
				}
				mgr.config.Debug("Broadcast on")
			} else {
				mgr.config.Debug("Sending message to stop the broadcast")
				err := mgr.stopBroadcast()
				if err != nil {
					mgr.config.Error("Error when stopping broadcast: %v", err)
					continue
				}
				mgr.config.Debug("Broadcast off")
			}

		case _ = <-mgr.scheduleChan:
			mgr.config.Notice("Schedule manager terminating")
			return
		}
	}
}
