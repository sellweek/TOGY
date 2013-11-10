package manager

import (
	"fmt"
	"github.com/sellweek/TOGY/control"
	"github.com/sellweek/TOGY/util"
	"os"
)

func startRotator(mgr *Manager) chan<- bool {
	exitChan := make(chan bool)
	go func() {
		mgr.config.Debug("Rotator started with presentations: %v", mgr.activePresentations)
		if len(mgr.activePresentations) != 0 {
			for {
				for _, p := range mgr.activePresentations {
					select {
					case <-exitChan:
						mgr.config.Debug("Rotator exiting")
						return
					default:
						mgr.config.Debug("Starting presentation: %s", p)
						pth, err := getPresentation(fmt.Sprint(mgr.config.BroadcastDir, string(os.PathSeparator), p))
						if err != nil {
							mgr.config.Error("Rotator couldn't get presentation: %v", err)
							continue
						}
						presentation := control.NewPowerPoint(mgr.config.PowerPoint, pth)
						mgr.config.Notice("New presentation was created")
						err = presentation.Run()
						if err != nil {
							mgr.config.Error("Rotator couldn't start PowerPoint: %v", err)
							continue
						}
					}
				}
			}
		} else {
			mgr.config.Debug("Rotator running without any presentations")
			<-exitChan
			mgr.config.Debug("Rotator exiting")
			return
		}
	}()
	return exitChan
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
	if fn != "" {
		return dir + string(os.PathSeparator) + fn, nil
	} else {
		return "", fmt.Errorf("Couldn't find PowerPoint file in folder %s.", dir)
	}

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
