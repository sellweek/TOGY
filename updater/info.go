package updater

import (
	"TOGY/config"
	"encoding/json"
	"io"
	"net/http"
)

type info struct {
	Broadcast bool
	FileType  bool
	Config    bool
}

func GetInfo(c config.Config) (i info, err error) {
	r, err := downloadInfo("/update?client=" + c.Name)
	defer r.Close()
	if err != nil {
		return
	}

	return parseInfo(r)
}

func downloadInfo(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func parseInfo(r io.Reader) (i info, err error) {
	d := json.NewDecoder(r)
	err = d.Decode(i)
	return
}
