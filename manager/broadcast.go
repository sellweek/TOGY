package manager

import (
	"TOGY/control"
	"os"
	"strings"
)

func broadcastManager(mgr *Manager) {
	var presentation *control.PowerPointBroadcast
	presentation = nil
	for msg := range mgr.broadcastChan {
		switch msg {
		case startBroadcast:
			if presentation == nil {
				pth, err := getPresentation(mgr.config.BroadcastDir)
				if err != nil {
					mgr.broadcastErr <- err
					continue
				}
				presentation = control.NewPowerPoint(mgr.config.PowerPoint, pth)
				mgr.config.Log.Println("New presentation was created")
			}
			err := presentation.Start()
			if err != nil {
				mgr.broadcastErr <- err
				continue
			}
			mgr.config.Log.Println("Turning screen on")
			err = control.TurnScreenOn()
			if err != nil {
				mgr.broadcastErr <- err
				continue
			}

		case stopBroadcast:
			mgr.config.Log.Println("Turning screen off")
			err := control.TurnScreenOff()
			if err != nil {
				mgr.broadcastErr <- err
				continue
			}

			if presentation == nil {
				continue
			}
			err = presentation.Kill()
			if err != nil {
				mgr.broadcastErr <- err
				continue
			}

			presentation = nil
			mgr.config.Log.Println("The presentation was stopped")

		case block:
			mgr.config.Log.Println("Broadcast manager blocked.")
			for m := range mgr.broadcastChan {
				if m == unblock {
					mgr.config.Log.Println("Broadcast manager unblocked.")
					break
				}
			}
		}
	}
	mgr.config.Log.Println("Broadcast manager terminating")
}

func getPresentation(dir string) (string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	files, err := f.Readdirnames(0)
	if err != nil {
		return "", err
	}

	fn := getFileWithType("pptx", files)
	if fn == "" {
		fn = getFileWithType("ppt", files)
	}
	return dir + string(os.PathSeparator) + fn, nil
}

func getFileWithType(ft string, fns []string) string {
	for _, fn := range fns {
		if getFileType(fn) == ft {
			return fn
		}
	}
	return ""
}

func getFileType(fn string) string {
	parts := strings.Split(fn, ".")
	return parts[len(parts)-1]
}
