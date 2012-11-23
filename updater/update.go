package updater

import (
	"TOGY/config"
	"TOGY/control"
	"TOGY/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const broadcastPath = "activeBroadcast"

type updateInfo struct {
	Broadcast bool
	FileType  string
	Config    bool
}

func getUpdateInfo(url, name string) (upd updateInfo, err error) {
	resp, err := http.Get(url + "?client=" + name)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(respBody, &upd)
	return
}

//Updates presentation and configuration files if needed.
//Its first return value is new broadcast process, if the
//broadcast was updated.
//Second return value is true if configuration was updated, 
//because that makes restart of almost the whole program necessary.
func Update(c config.Config, runningBroadcast control.Broadcast) (b control.Broadcast, conf bool) {
	conf = false
	ui, err := getUpdateInfo(c.UpdateURL, c.Name)
	c.Log.Println("Got update info:", ui, err)
	if err != nil {
		c.Log.Println("Could not get update info: ", err)
		return
	}

	if ui.Broadcast {
		c.Log.Println("Downloading new broadcast.")
		b, err = getBroadcast(ui, c, runningBroadcast)
		if err != nil {
			c.Log.Println("Error when updating broadcast: ", err)
			return nil, false
		} else {
			c.Log.Println("Downloaded new broadcast.")
		}
	}

	if ui.Config {
		c.Log.Println("Downloading new config")
		err = getConfig(c)
		if err != nil {
			c.Log.Println("Error when downloading new config: ", err)
			conf = false
		} else {
			conf = true
			c.Log.Println("Downloaded new config.")
		}
	}
	return
}

func getConfig(c config.Config) error {
	err := downloadFile(c.UpdateURL+"/config?client="+c.Name, c.CentralPath)
	if err != nil {
		return fmt.Errorf("Could not download or save new configuration: %v", err)
	}

	_, err = http.Get(c.UpdateURL + "/gotConfig" + "?client=" + c.Name)
	if err != nil {
		c.Log.Println("Could not announce succesful download of configuration:", err)
	}
	return nil
}

func getBroadcast(ui updateInfo, c config.Config, runningBroadcast control.Broadcast) (b control.Broadcast, err error) {
	b = nil
	err = downloadFile(c.UpdateURL+"/download?client="+c.Name, "newBroadcast."+ui.FileType)
	if err != nil {
		err = fmt.Errorf("Could not download or save presentation: %v", err)
		return
	}

	//It should be safe to call Kill(), becuase runningBroadcast should not be a nil pointer
	//because we are starting updater goroutine after 30 second delay
	//and Kill() is a no-op if the broadcast has already been killed.
	runningBroadcast.Kill()
	//We have to wait a bit to ensure that handler application
	//will not block the removal of the file.
	util.Sleep(1)
	//We have to remove the current broadcast, because its file type
	//could be different from the one we will download.
	err = os.Remove(runningBroadcast.Path())
	if err != nil {
		err = fmt.Errorf("Could not remove current broadcast: %v", err)
		return
	}

	err = os.Rename("newBroadcast."+ui.FileType, broadcastPath+"."+ui.FileType)
	if err != nil {
		err = fmt.Errorf("Could not move presentation: %v", err)
		return
	}

	util.Sleep(1)
	currPath, err := GetCurrentBroadcast(".")
	if err != nil {
		err = fmt.Errorf("Could not find current broadcast: %v", err)
		return
	}

	b = control.NewPowerPoint(c.PowerPoint, currPath)

	if err = b.Start(); err != nil {
		fmt.Errorf("Could not start broadcast: %v", err)
		return
	}

	_, err = http.Get(c.UpdateURL + "/downloadComplete" + "?client=" + c.Name)
	if err != nil {
		c.Log.Println("Could not announce succesful download of broadcast:", err)
	}
	return
}

//Downloads file from specified URL to specified path.
func downloadFile(src, dest string) (err error) {
	file, err := os.Create(dest)
	if err != nil {
		return
	}

	resp, err := http.Get(src)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return
	}

	err = file.Close()
	if err != nil {
		return
	}

	return
}
