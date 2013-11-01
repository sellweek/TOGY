package manager

import (
	"fmt"
	"github.com/sellweek/TOGY/updater"
	"os"
	"strconv"
	"time"
)

//updateManager receives messages from a Ticker,
//checking server for update info and downloading
//updated config or broadcast.
func updateManager(mgr *Manager, t <-chan time.Time) {
	/*for {
		select {
		case _ = <-t:
			mgr.config.Notice("Getting update info")
			i, err := updater.GetInfo(mgr.config)
			if err != nil {
				mgr.config.Error("Error while downloading info: %v", err)
				continue
			}
			if i.Broadcast {
				mgr.config.Notice("Updating broadcast")
				err = updateBroadcast(mgr, i.FileType)
				if err != nil {
					mgr.config.Error("Error when updating broadcast: %v", err)
				}
				mgr.config.Notice("Broadcast updated")
			}

			//If a new config is downloaded, a reload signal is sent
			//and updateManager terminates.
			if i.Config {
				mgr.config.Notice("Updating config")
				err = updateConfig(mgr)
				if err != nil {
					mgr.config.Error("Error when updating config: %v", err)
				}
				mgr.reloadSignal <- true
				mgr.config.Notice("Config updated, restarting manager.")
				return
			}
		}
	}*/
}

//updateConfig downloads a new config from server and announces its
//succesful download back.
func updateConfig(mgr *Manager) error {
	err := updater.DownloadConfig(mgr.config, mgr.config.CentralPath)
	if err != nil {
		return fmt.Errorf("Error while updating config: %v", err)
	}
	mgr.config.Notice("Config succesfully downloaded, announcing.")
	err = updater.AnnounceConfig(mgr.config)
	if err != nil {
		return fmt.Errorf("Error while announcing the download of config: %v", err)
	}
	return nil
}

//updateBroadcast downloads a new broadcast from server,
//announces its succesful download and moves it into place.
func updateBroadcast(mgr *Manager, ft string) error {
	/*path, err := makeTempDir()
	if err != nil {
		return fmt.Errorf("Error while creating temporary directory: %v", err)
	}
	err = updater.DownloadBroadcast(mgr.config, ft, path)
	if err != nil {
		return fmt.Errorf("Error while downloading new broadcast: %v", err)
	}
	mgr.config.Notice("Broadcast successfully downloaded, announcing.")

	err = updater.AnnounceBroadcast(mgr.config)
	if err != nil {
		return fmt.Errorf("Error while announcing the download of broadcast: %v", err)
	}

	err = mgr.stopBroadcast()
	if err != nil {
		return fmt.Errorf("Error while stopping broadcast for update: %v", err)
	}

	//scheduleManager has to be blocked from interfering
	//with us moving the broadcast into place.
	mgr.block()
	defer mgr.unblock()

	mgr.config.Notice("Moving new broadcast into place")
	err = deleteAll(mgr.config.BroadcastDir)
	if err != nil {
		return fmt.Errorf("Error while deleting old broadcast: %v", err)
	}

	err = moveFiles(path, mgr.config.BroadcastDir)
	if err != nil {
		return fmt.Errorf("Error while moving new broadcast into place: %v", err)
	}

	err = mgr.startBroadcast()
	if err != nil {
		return fmt.Errorf("Error while starting new broadcast: %v", err)
	}

	_ = deleteAll(path)*/
	return nil
}

//makeTempDir creates a directory in the temporary directory
//returned by os.TempDir() and returns its path.
func makeTempDir() (path string, err error) {
	path = os.TempDir() + string(os.PathSeparator) + "broadcast-download-" + strconv.Itoa(int(time.Now().Unix()))
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return
	}
	return
}

//moveFiles moves all the files from one directory into another.
func moveFiles(src, dest string) (err error) {
	sf, err := os.Open(src)
	if err != nil {
		return
	}

	fns, err := sf.Readdirnames(0)
	if err != nil {
		return
	}

	for _, name := range fns {
		path := src + string(os.PathSeparator) + name
		err = os.Rename(path, dest+string(os.PathSeparator)+name)
		if err != nil {
			return
		}
	}
	return
}

//deleteAll deletes all the files in given directory.
func deleteAll(dir string) (err error) {
	d, err := os.Open(dir)
	if err != nil {
		return
	}

	fns, err := d.Readdirnames(0)
	if err != nil {
		return
	}

	for _, f := range fns {
		path := dir + string(os.PathSeparator) + f
		err = os.Remove(path)
		if err != nil {
			return
		}
	}
	return nil
}
