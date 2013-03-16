package updater

import (
	"archive/zip"
	"github.com/sellweek/TOGY/config"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

//DownloadConfig ownloads a new centralConfig.json from the server
//into the given path.
func DownloadConfig(c *config.Config, destFile string) error {
	return downloadFile(c.UpdateURL+"/config/download?client="+c.Name, destFile)
}

//DownloadBroadcast downloads a new broadcast from the server into a 
//given directory, unzipping it, if it has .zip extension.
func DownloadBroadcast(c *config.Config, ft string, destDir string) (err error) {
	srcUrl := c.UpdateURL + "/presentation/active/download?client=" + c.Name

	if ft != "zip" {
		err = downloadFile(srcUrl, destDir+string(os.PathSeparator)+"broadcast."+ft)
		return
	}

	tempFileName := os.TempDir() + string(os.PathSeparator) + "unzip-" + strconv.Itoa(int(time.Now().Unix())) + ".zip"

	err = downloadFile(srcUrl, tempFileName)
	if err != nil {
		return
	}

	err = unzip(destDir, tempFileName)

	return
}

//ColdStart downloads central config and the newest broadcast
//into folders specified in config, without announcing
//their downloads.
func ColdStart(c *config.Config) (err error) {
	ui, err := GetInfo(c)
	if err != nil {
		return err
	}
	err = DownloadConfig(c, c.CentralPath)
	if err != nil {
		return err
	}
	err = DownloadBroadcast(c, ui.FileType, c.BroadcastDir)
	return
}

//Unzip unzips a file into given folder.
//
//WARNING: the unzipping is not recursive therefore it doesn't support
//zip files with folders.
func unzip(dirname, fn string) (err error) {
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
		df, err = os.Create(dirname + string(os.PathSeparator) + sf.Name)
		if err != nil {
			return
		}
		defer df.Close()
		_, err = io.Copy(df, fr)
		if err != nil {
			return
		}
	}
	return nil
}

//Downloads a file from given URL into given path using http.Get
func downloadFile(src, dest string) (err error) {
	resp, err := http.Get(src)
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return
}
