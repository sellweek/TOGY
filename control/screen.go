package control

import (
	"os/exec"
)

//TurnScreenOn turns computer's screen on
func TurnScreenOn() error {
	return setScreenState("on")
}

//TurnScreenOff turns computer's screen off
func TurnScreenOff() error {
	return setScreenState("off")
}

func setScreenState(state string) error {
	cmd := exec.Command("nircmd", "monitor", state)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
