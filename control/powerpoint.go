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

//Start starts PowerPoint in presentation mode.
func (b *PowerPointBroadcast) Start() error {
	cmd := exec.Command(b.powerPoint, "/s", b.path)
	err := cmd.Start()
	if err != nil {
		return err
	}
	b.cmd = cmd
	return nil
}

//Kill ends Terminate signal to PowerPoint.
func (b *PowerPointBroadcast) Kill() error {
	if b.cmd == nil {
		return nil
	}
	err := b.cmd.Process.Kill()
	if err != nil {
		return err
	}
	b.cmd = nil
	return nil
}

func (b PowerPointBroadcast) Status() bool {
	if b.cmd == nil {
		return false
	}
	return true
}

func (b PowerPointBroadcast) Path() string {
	return b.path
}

func NewPowerPoint(ppExe, presentation string) *PowerPointBroadcast {
	return &PowerPointBroadcast{path: presentation, powerPoint: ppExe}
}
