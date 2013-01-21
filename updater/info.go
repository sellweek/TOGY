package updater

import (
	"TOGY/config"
	"encoding/json"
	"io"
	"net/http"
)

type Info struct {
	Broadcast bool
	FileType  string
	Config    bool
}

func GetInfo(c *config.Config) (i Info, err error) {
	r, err := downloadInfo(c.UpdateURL + "/update?client=" + c.Name)
	if err != nil {
		return
	}
	defer r.Close()

	return parseInfo(r)
}

func downloadInfo(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func parseInfo(r io.Reader) (i Info, err error) {
	d := json.NewDecoder(r)
	err = d.Decode(&i)
	return
}

func AnnounceBroadcast(c *config.Config) error {
	url := c.UpdateURL + "/presentation/active/downloadComplete?client=" + c.Name
	_, err := http.Get(url)
	return err
}

func AnnounceConfig(c *config.Config) error {
	url := c.UpdateURL + "/cofig/downloadComplete?client=" + c.Name
	_, err := http.Get(url)
	return err
}
