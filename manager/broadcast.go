package manager

import (
	"fmt"
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
	var (
		errChan  <-chan error
		stopChan chan<- bool
	)
	for msg := range mgr.broadcastChan {
		mgr.config.Debug("Recieved message: %v", msg)
		switch msg {
		//When broadcast manager receives a message
		//telling it to turn the broadcast on,
		//it starts the handler application
		//and turns the screen on.
		case startBroadcast:
			if stopChan == nil {
				stopChan, errChan = presentationRotator(mgr)
				mgr.config.Debug("Turning screen on")
				err := control.TurnScreenOn()
				if err != nil {
					mgr.broadcastErr <- err
					continue
				}
				mgr.broadcastErr <- nil
			} else {
				select {
				case err := <-errChan:
					mgr.broadcastErr <- err
					continue
				default:
					mgr.broadcastErr <- nil
					continue
				}
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

			if stopChan == nil {
				mgr.broadcastErr <- nil
				continue
			}
			stopChan <- true

			stopChan = nil
			errChan = nil
			mgr.config.Notice("The presentation was stopped")
			mgr.broadcastErr <- nil

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

func presentationRotator(mgr *Manager) (chan<- bool, <-chan error) {
	exitChan := make(chan bool)
	errChan := make(chan error)
	go func() {
		mgr.config.Debug("Rotator started with presentations: %v", mgr.currentPresentations)
		for {
			for _, p := range mgr.currentPresentations {
				mgr.config.Debug("Starting presentation: %s", p)
				select {
				case <-exitChan:
					mgr.config.Debug("Rotator exiting")
					return
				default:
					pth, err := getPresentation(fmt.Sprint(mgr.config.BroadcastDir, string(os.PathSeparator), p))
					if err != nil {
						mgr.config.Error("Rotator couldn't get presentation: %v", err)
						errChan <- err
						continue
					}
					presentation := control.NewPowerPoint(mgr.config.PowerPoint, pth)
					mgr.config.Notice("New presentation was created")
					err = presentation.Run()
					if err != nil {
						mgr.config.Error("Rotator couldn't start PowerPoint: %v", err)
						errChan <- err
						continue
					}
				}
			}
		}
	}()
	return exitChan, errChan
}
