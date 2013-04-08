package manager

import (
	"github.com/sellweek/TOGY/control"
	"github.com/sellweek/TOGY/util"
	"os"
)

//broadcastManager takes care of
//turning the handler program and screen on and off.
//It receives messages from the schedule manager
//and starts and stops broadcast according to them.
//All errors that occur are sent back on 
//mgr.broadcastErr channel.
//It can also be blocked, which makes it
//wait for unblock message, throwing away all the
//other messages.
func broadcastManager(mgr *Manager) {
	var presentation *control.PowerPointBroadcast
	presentation = nil
	for msg := range mgr.broadcastChan {
		switch msg {
		//When broadcast manager receives a message
		//telling it to turn the broadcast on,
		//it starts the handler application
		//and turns the screen on. 
		case startBroadcast:
			if presentation == nil {
				pth, err := getPresentation(mgr.config.BroadcastDir)
				if err != nil {
					mgr.broadcastErr <- err
					continue
				}
				presentation = control.NewPowerPoint(mgr.config.PowerPoint, pth)
				mgr.config.Notice("New presentation was created")
			}
			err := presentation.Start()
			if err != nil {
				mgr.broadcastErr <- err
				continue
			}
			mgr.config.Debug("Turning screen on")
			err = control.TurnScreenOn()
			if err != nil {
				mgr.broadcastErr <- err
				continue
			}

		//When broadcast manager receives a message
		//telling it to stop the broadcast,
		//it terminates the handler application
		//and turns the screen off.
		case stopBroadcast:
			mgr.config.Debug("Turning screen off")
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
			mgr.config.Notice("The presentation was stopped")

		//When broadcast manager receives a message
		//telling it to block, it will throw away all
		//messages received, unless they tell it to unblock.
		case block:
			mgr.config.Info("Broadcast manager blocked.")
			for m := range mgr.broadcastChan {
				if m == unblock {
					mgr.config.Info("Broadcast manager unblocked.")
					break
				}
			}
		}
	}
	mgr.config.Notice("Broadcast manager terminating")
}

//getPresentation searches the given directory
//for a file with "ppt" or "pptx" extension
//and returns its name.
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

//Searches a list of file names for the file with
//a given extension and returns its name.
func getFileWithType(ft string, fns []string) string {
	for _, fn := range fns {
		if util.GetFileType(fn) == ft {
			return fn
		}
	}
	return ""
}
