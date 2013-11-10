package manager

import (
	"github.com/sellweek/TOGY/control"
	"time"
)

func broadcastManager(m *Manager, st <-chan time.Time) {
	m.config.Debug("Broadcast manager started")
	var rotator chan<- bool
	for {
		select {
		case <-st:
			if m.config.Broadcast() {
				err := control.TurnScreenOn()
				if err != nil {
					m.config.Error("Error while turning screen on: %v", err)
				}
				if rotator == nil {
					rotator = startRotator(m)
				}
			} else {
				if rotator != nil {
					rotator <- true
					rotator = nil
				}
				err := control.TurnScreenOff()
				if err != nil {
					m.config.Error("Error while turning screen off: %v", err)
				}
			}
		case <-m.kbChan:
			if rotator != nil {
				rotator <- true
				rotator = nil
				m.config.Debug("Rotator exited")
			}
			err := control.TurnScreenOff()
			if err != nil {
				m.config.Error("Error while turning screen off: %v", err)
			}
			m.broadcastKilled <- true
			m.config.Debug("Broadcast exited")
			return
		}
	}
}
