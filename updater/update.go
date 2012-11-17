package updater

import (
	"TOGY/config"
	"TOGY/control"
	"TOGY/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
)

const broadcastPath = "activeBroadcast"

type updateInfo struct {
	Broadcast bool
	FileType string
	Config bool
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
		b = getBroadcast(ui, c, runningBroadcast)
	}

	if ui.Config {
		getConfig(c)
		conf = true
	}
	return
}

func getConfig(c config.Config) {
	c.Log.Println("Downloading new config")
	err := downloadFile(c.UpdateURL+"/config?client=" + c.Name, c.CentralPath)
	if err != nil {
		c.Log.Println("Could not download or save new configuration: ", err)
	}

	_, err = http.Get(c.UpdateURL + "/gotConfig" + "?client=" + c.Name)
	if err != nil {
		c.Log.Println("Could not announce succesful download of configuration:", err)
	}

}


func getBroadcast(ui updateInfo, c config.Config, runningBroadcast control.Broadcast) control.Broadcast {

	c.Log.Println("Downloading new version")
	err := downloadFile(c.UpdateURL + "/download?client=" + c.Name, "newBroadcast."+ui.FileType)
	if err != nil {
		c.Log.Println("Could not download or save presentation:", err)
		return nil
	}
	if runningBroadcast != nil {
		runningBroadcast.Kill()
		util.Sleep(1)
		//We have to remove the current broadcast, because its file type
		//could be different from the one we will download.
		err = os.Remove(runningBroadcast.Path())
	} else {
		cb, err := GetCurrentBroadcast(".")
		if err != nil {
			c.Log.Println("Could not find current broadcast:", err)
		}
		err = os.Remove(cb)
	}

	if err != nil {
		c.Log.Println("Could not remove current broadcast: ", err)
	}


	err = os.Rename("newBroadcast."+ui.FileType, broadcastPath+"."+ui.FileType)
	if err != nil {
		c.Log.Println("Could not move presentation:", err)
		return nil
	}
	
	util.Sleep(1)
	currPath, err := GetCurrentBroadcast(".")
	if err != nil {
		c.Log.Println("Could not find current broadcast:", err)
	}

	var currCast control.Broadcast = control.NewPowerPoint(c.PowerPoint, currPath)

	if err := currCast.Start(); err != nil {
		c.Log.Println("Could not start broadcast:", err)
		return nil
	}

	_, err = http.Get(c.UpdateURL + "/downloadComplete" + "?client=" + c.Name)
	if err != nil {
		c.Log.Println("Could not announce succesful download of broadcast:", err)
	}
	return currCast
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