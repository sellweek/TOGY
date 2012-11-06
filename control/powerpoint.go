package control

import (
	"TOGY/util"
	"os/exec"
)

// Starts PowerPoint in presentation mode.
func StartPresentation(ppExe, path string) error {
	err := exec.Command(ppExe, "/s", path).Start()
	if err != nil {
		return err
	}
	return nil
}

//Sends Terminate signal to PowerPoint.
func KillPresentation() {
	exec.Command("taskkill", "/IM", "POWERPNT.exe").Run()
}

//Kills powerpoint and loads presentation at p.
func ReloadPresentation(ppExe, p string) {
	KillPresentation()
	util.Sleep(1)
	StartPresentation(ppExe, p)
}
