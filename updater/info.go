package updater

import (
	"encoding/json"
	"github.com/sellweek/TOGY/config"
	"io"
	"net/http"
)

//Info is a struct used to represent information about
//updated broadcast and config on the server.
type Info struct {
	//Broadcast informs whether there is
	//a newer broadcast than the one client has.
	Broadcast bool

	//FileType gives the type of the broadcast file.
	FileType string

	//Config informs whether there is a newer version
	//of centralConfig available.
	Config bool
}

//GetInfo returns an Info struct with current
//information from server.
func GetInfo(c *config.Config) (i Info, err error) {
	r, err := downloadInfo(c.UpdateURL + "/update?client=" + c.Name)
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

//AnnounceBroadcast announces the completition of download 
//of the broadcast to the server.
func AnnounceBroadcast(c *config.Config) error {
	url := c.UpdateURL + "/presentation/active/downloadComplete?client=" + c.Name
	_, err := http.Get(url)
	return err
}

//AnnounceConfig announces the completition of download 
//of the centralConfig to the server.
func AnnounceConfig(c *config.Config) error {
	url := c.UpdateURL + "/config/downloadComplete?client=" + c.Name
	_, err := http.Get(url)
	return err
}
