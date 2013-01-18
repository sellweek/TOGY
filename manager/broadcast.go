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

			}
			err := presentation.Start()
			if err != nil {
				mgr.broadcastErr <- err
				continue
			}
			err = control.TurnScreenOn()
			if err != nil {
				mgr.broadcastErr <- err
				continue
			}

		case stopBroadcast:
			if presentation == nil {
				continue
			}
			err := presentation.Kill()
			if err != nil {
				mgr.broadcastErr <- err
				continue
			}

			err = control.TurnScreenOff()
			if err != nil {
				mgr.broadcastErr <- err
				continue
			}
			presentation = nil

		case block:
			for m := range mgr.broadcastChan {
				if m == unblock {
					break
				}
			}
		}
	}
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
