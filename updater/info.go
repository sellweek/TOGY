package updater

import (
	"encoding/json"
	"fmt"
	"github.com/sellweek/TOGY/config"
	"io"
	"net/http"
)

type BroadcastInfo struct {
	Key      string
	FileType string
}

type Info struct {
	Broadcasts []BroadcastInfo
	Config     int64
}

//GetInfo returns an Info struct with current
//information from server.
func GetInfo(c *config.Config) (i Info, err error) {
	r, err := downloadInfo(c.UpdateURL + "/status?client=" + c.Name)
	if err != nil {
		return
	}
	defer r.Close()
	return parseInfo(r)
}

//downloadInfo downloads the update JSON from the server.
func downloadInfo(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

//parseInfo parses the JSON info from the server.
func parseInfo(r io.Reader) (i Info, err error) {
	d := json.NewDecoder(r)
	err = d.Decode(&i)
	return
}

//AnnounceBroadcast announces to the server that
//the broadcast has been succesfuly activated.
func AnnounceBroadcast(c *config.Config, key string) error {
	url := fmt.Sprint(c.UpdateURL, "/presentation/", key, "/downloadComplete?client=", c.Name)
	_, err := http.Get(url)
	return err
}

//AnnounceBroadcast announces to the server that
//the config has been succesfuly downloaded.
func AnnounceConfig(c *config.Config) error {
	url := c.UpdateURL + "/config/downloadComplete?client=" + c.Name
	_, err := http.Get(url)
	return err
}
