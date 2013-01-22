package manager

import (
	"fmt"
	"github.com/sellweek/TOGY/updater"
	"os"
	"strconv"
	"time"
)

func updateManager(mgr *Manager, t <-chan time.Time) {
	for {
		select {
		case _ = <-t:
			mgr.config.Log.Println("Getting update info")
			i, err := updater.GetInfo(mgr.config)
			if err != nil {
				mgr.config.Log.Println("Error while downloading info: ", err)
				continue
			}
			if i.Broadcast {
				mgr.config.Log.Println("Updating broadcast")
				err = updateBroadcast(mgr, i.FileType)
				if err != nil {
					mgr.config.Log.Println(err)
				}
				mgr.config.Log.Println("Broadcast updated")
			}

			if i.Config {
				mgr.config.Log.Println("Updating config")
				err = updateConfig(mgr)
				if err != nil {
					mgr.config.Log.Println(err)
				}
				mgr.reloadSignal <- true
				mgr.config.Log.Println("Config updated, restarting manager.")
				return
			}
		}
	}
}

func updateConfig(mgr *Manager) error {
	err := updater.DownloadConfig(mgr.config, mgr.config.CentralPath)
	if err != nil {
		return fmt.Errorf("Error while updating config: %v", err)
	}
	mgr.config.Log.Println("Config succesfully downloaded, announcing.")
	err = updater.AnnounceConfig(mgr.config)
	if err != nil {
		return fmt.Errorf("Error while announcing the download of config: %v", err)
	}
	return nil
}

func updateBroadcast(mgr *Manager, ft string) error {
	path, err := makeTempDir()
	if err != nil {
		return fmt.Errorf("Error while creating temporary directory: %v", err)
	}
	err = updater.DownloadBroadcast(mgr.config, ft, path)
	if err != nil {
		return fmt.Errorf("Error while downloading new broadcast: %v", err)
	}
	mgr.config.Log.Println("Broadcast successfully downloaded, announcing.")

	err = updater.AnnounceBroadcast(mgr.config)
	if err != nil {
		return fmt.Errorf("Error while announcing the download of broadcast: %v", err)
	}

	err = mgr.stopBroadcast()
	if err != nil {
		return fmt.Errorf("Error while stopping broadcast for update: %v", err)
	}

	mgr.block()
	defer mgr.unblock()

	mgr.config.Log.Println("Moving new broadcast into place")
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

	_ = deleteAll(path)
	return nil
}

func makeTempDir() (path string, err error) {
	path = os.TempDir() + string(os.PathSeparator) + "broadcast-download-" + strconv.Itoa(int(time.Now().Unix()))
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return
	}
	return
}

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
