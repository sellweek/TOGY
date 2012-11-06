package updater

import (
	"os"
	"strings"
	"errors"
)

func GetCurrentBroadcast(dir string) (string, error) {
	folder, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	files, err := folder.Readdir(0)
	if err != nil {
		return "", err
	}
	for _, fi := range files {
		fl := strings.Split(fi.Name(), ".")
		if fl[0] == broadcastPath {
			return fi.Name(), nil
		}
	}
	return "", errors.New("No presentation found.")
}