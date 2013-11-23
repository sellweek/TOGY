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
	for {
		select {
		case _ = <-t:
			mgr.config.Notice("Getting update info")
			i, err := updater.GetInfo(mgr.config)
			if err != nil {
				mgr.config.Error("Error while downloading info: %v", err)
				continue
			}

			err = updateBroadcasts(mgr, i)
			if err != nil {
				mgr.config.Error("Error while updating broadcasts: %v", err)
				continue
			}

			//If a new config is downloaded, a reload signal is sent
			//and updateManager terminates.
			if i.Config > mgr.config.Timestamp {
				mgr.config.Notice("Updating config")
				err = updateConfig(mgr)
				if err != nil {
					mgr.config.Error("Error when updating config: %v", err)
				}
				mgr.reload <- true
				mgr.config.Notice("Config updated, restarting manager.")
				return
			}
		}
	}
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
func updateBroadcasts(mgr *Manager, inf updater.Info) (err error) {
	addedBroadcasts := make([]updater.BroadcastInfo, 0)
	for _, p := range inf.Broadcasts {
		if !stringInSlice(p.Key, mgr.activePresentations) {
			addedBroadcasts = append(addedBroadcasts, p)
		}
	}
	mgr.config.Debug("Added broadcasts: %v", addedBroadcasts)

	removedBroadcasts := make([]string, 0)
	bKeys := make([]string, len(inf.Broadcasts))
	for i, p := range inf.Broadcasts {
		bKeys[i] = p.Key
	}

	for _, p := range mgr.activePresentations {
		if !stringInSlice(p, bKeys) {
			removedBroadcasts = append(removedBroadcasts, p)
		}
	}
	mgr.config.Debug("Removed broadcasts: %v", removedBroadcasts)

	wasUpdated := len(addedBroadcasts) != 0 || len(removedBroadcasts) != 0

	if wasUpdated {
		mgr.killBroadcast()
		defer mgr.startBroadcast()
	}

	if len(addedBroadcasts) != 0 {
		path, err := makeTempDir()
		if err != nil {
			mgr.config.Error("Error while creating temporary directory: %v", err)
		}
		err = updater.DownloadBroadcasts(mgr.config, addedBroadcasts, path)
		if err != nil {
			mgr.config.Error("Error while downloading new broadcasts: %v", err)
		}
		mgr.config.Notice("Broadcasts successfully downloaded, moving into broadcast folder.")

		err = moveFiles(path, mgr.config.BroadcastDir)
		if err != nil {
			mgr.config.Error("Error while moving new broadcast into place: %v", err)
		}

		_ = deleteAll(path)
	}

	if len(removedBroadcasts) != 0 {
		for _, key := range removedBroadcasts {
			rmPath := fmt.Sprint(mgr.config.BroadcastDir, string(os.PathSeparator), key)
			mgr.config.Debug("Removing directory: %s", rmPath)
			err = os.RemoveAll(rmPath)
			if err != nil {
				mgr.config.Error("Error while removing deactivated broadcast %v: %v", key, err)
			}
		}
	}

	err = nil

	if wasUpdated {
		mgr.activePresentations, err = getBroadcastDirs(mgr.config)
		if err != nil {
			return fmt.Errorf("Couldn't get current broadcasts: %v", err)
		}
	}

	for _, p := range addedBroadcasts {
		mgr.config.Debug("Announcing broadcast %v", p.Key)
		err = updater.AnnounceBroadcast(mgr.config, p.Key)
		if err != nil {
			mgr.config.Error("Error while announcing activation of broadcast %v: %v", p.Key, err)
		}
	}
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
