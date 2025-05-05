package paths

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/smoothie-go/smgo-manager/hlog"
)

func userDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		dir := os.Getenv("LocalAppData")
		if dir == "" {
			dir = os.Getenv("AppData")
		}
		return dir, nil
	case "darwin":
		return filepath.Join(home, "Library", "Application Support"), nil
	case "linux":
		dir := os.Getenv("XDG_DATA_HOME")
		if dir == "" {
			dir = filepath.Join(home, ".local", "share")
		}
		return dir, nil
	default:
		return "", errors.New("Operating systems other than Linux, macOS, and Windows are not supported")
	}
}

func getManagerDirectory() string {
	data_dir, err := userDataDir()
	if err != nil {
		hlog.Fatal(err.Error())
	}
	path := filepath.Join(data_dir, "smgo-manager")
	d, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			hlog.Fatal(err.Error())
		}
	} else if !d.IsDir() {
		hlog.Fatal(path + " Isn't a directory")
	}

	return path
}

func manDirJoin(name string) string {
	man_dir := getManagerDirectory()

	path := filepath.Join(man_dir, name)
	d, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			hlog.Fatal(err.Error())
		}
	} else if !d.IsDir() {
		hlog.Fatal(path + " Isn't a directory")
	}

	return path
}

func GetVersionsDirectory() string {
	return manDirJoin("Versions")
}

func GetDownloadsDirectory() string {
	return manDirJoin("Downloads")
}

func GetTempDirectory() string {
	return manDirJoin("Temp")
}
