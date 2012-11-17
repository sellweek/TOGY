package control

import (
	"os/exec"
)

type PowerPointBroadcast struct {
	//Path to the presentation
	path string
	//Path to PowerPoint executable
	powerPoint string
	//Current running PowerPoint instance
	cmd *exec.Cmd
}

// Starts PowerPoint in presentation mode.
func(b *PowerPointBroadcast) Start() error {
	cmd := exec.Command(b.powerPoint, "/s", b.path)
	err := cmd.Start()
	if err != nil {
		return err
	}
	b.cmd = cmd
	return nil
}

//Sends Terminate signal to PowerPoint.
func(b *PowerPointBroadcast) Kill() error {
	err := b.cmd.Process.Kill()
	if err != nil {
		return err
	}
	b.cmd = nil
	return nil
}

func(b PowerPointBroadcast) Status() bool {
	if b.cmd == nil {
		return false
	}
	return true
}

func NewPowerPoint(ppExe, presentation string) (*PowerPointBroadcast) {
	return &PowerPointBroadcast{path: presentation, powerPoint: ppExe}
}