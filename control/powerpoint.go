package control

import (
	"os/exec"
)

//PowerPointBroadcast represents a broadcast
//int ppt or pptx format.
type PowerPointBroadcast struct {
	//Path to the presentation
	path string
	//Path to PowerPoint executable
	powerPoint string
	//Current running PowerPoint instance
	cmd *exec.Cmd
}

//Run starts PowerPoint in presentation mode
//and waits until it terminates.
func (b *PowerPointBroadcast) Run() error {
	if b.Status() {
		return nil
	}
	cmd := exec.Command(b.powerPoint, "/s", b.path)
	b.cmd = cmd
	err := cmd.Run()
	if err != nil {
		return err
	}
	b.cmd = nil
	return nil
}

//Kill ends Terminate signal to PowerPoint.
func (b *PowerPointBroadcast) Kill() error {
	if !b.Status() {
		return nil
	}
	err := b.cmd.Process.Kill()
	if err != nil {
		return err
	}
	b.cmd = nil
	return nil
}

//Status returns whether given broadcast is still
//broadcasting.
func (b PowerPointBroadcast) Status() bool {
	if b.cmd == nil {
		return false
	}
	return true
}

//Path returns the path to given broadcast's ppt or pptx file.
func (b PowerPointBroadcast) Path() string {
	return b.path
}

//NewPowerPoint returns a new PowerPointBroadcast
func NewPowerPoint(ppExe, presentation string) *PowerPointBroadcast {
	return &PowerPointBroadcast{path: presentation, powerPoint: ppExe}
}
