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
//Returns true if configuration was updated, because that
//makes restart of pretty much whole program necessary.
func Update(c config.Config) bool {
	ui, err := getUpdateInfo(c.UpdateURL, c.Name)
	c.Log.Println("Got update info:", ui, err)
	if err != nil {
		c.Log.Println("Could not get update info: ", err)
		return false
	}

	if ui.Broadcast {
		getBroadcast(ui, c)
	}

	if ui.Config {
		getConfig(c)
		return true
	}
	return false
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


func getBroadcast(ui updateInfo, c config.Config) {

	c.Log.Println("Downloading new version")
	err := downloadFile(c.UpdateURL + "/download?client=" + c.Name, "newBroadcast."+ui.FileType)
	if err != nil {
		c.Log.Println("Could not download or save presentation:", err)
		return
	}

	control.KillPresentation()

	current, err := GetCurrentBroadcast(".")
	if err != nil {
		c.Log.Println("Could not get current broadcast: ", err)
	}

	//We have to remove the current broadcast, because its file type
	//could be different from the one we will download.
	err = os.Remove(current)
	if err != nil {
		c.Log.Println("Could not remove current broadcast: ", err)
	}

	err = os.Rename("newBroadcast."+ui.FileType, broadcastPath+"."+ui.FileType)
	if err != nil {
		c.Log.Println("Could not move presentation:", err)
		return
	}
	
	util.Sleep(1)
	currCast, _ := GetCurrentBroadcast(".")

	control.StartPresentation(c.PowerPoint, currCast)

	_, err = http.Get(c.UpdateURL + "/downloadComplete" + "?client=" + c.Name)
	if err != nil {
		c.Log.Println("Could not announce succesful download of broadcast:", err)
	}
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