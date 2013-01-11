package updater

import (
	"TOGY/config"
	"archive/zip"
	"io"
	"net/http"
	"os"
	"time"
)

func DownloadConfig(c config.Config) (f *os.File, err error) {
	resp, err := http.Get(c.UpdateURL + "/config/download?client=" + c.Name)
	defer resp.Body.Close()
	f, err = createTempFile("config", "json", resp.Body)
	return
}

func DownloadBroadcast(c config.Config, ft string) (f *os.File, err error) {
	resp, err := http.Get(c.UpdateURL + "/presentation/active/download?client" + c.Name)
	defer resp.Body.Close()

	f, err = createTempFile("broadcast", ft, resp.Body)
	if err != nil {
		return
	}
	if ft != "zip" {
		return
	}

	defer f.Close()

	dirName := os.TempDir() + "/broadcast-zip-" + time.Now().String()
	err = os.Mkdir(dirName, os.ModePerm)
	if err != nil {
		return
	}

	fi, err := f.Stat()
	if err != nil {
		return
	}

	err = unzip(dirName, fi.Name())
	if err != nil {
		return
	}
	return os.Open(dirName)
}

func createTempFile(prefix, fileType string, data io.Reader) (f *os.File, err error) {
	name := os.TempDir() + prefix + "-" + time.Now().String() + "." + fileType
	f, err = os.Create(name)
	if err != nil {
		return
	}
	_, err = io.Copy(f, data)
	return
}

func unzip(dirname string, fn string) (err error) {
	r, err := zip.OpenReader(fn)
	if err != nil {
		return
	}

	for _, sf := range r.File {
		var fr io.ReadCloser
		fr, err = sf.Open()
		if err != nil {
			return
		}
		var df *os.File
		df, err = os.Create(dirname + "/" + sf.Name)
		if err != nil {
			return
		}
		_, err = io.Copy(df, fr)
		if err != nil {
			return
		}
	}
	return nil
}
