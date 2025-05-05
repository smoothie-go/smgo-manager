package download

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"github.com/smoothie-go/smgo-manager/hlog"
)

func Download(url, output string) {
	resp, err := http.Get(url)
	if err != nil {
		hlog.Fatal(err.Error())
	}
	defer resp.Body.Close()

	clen := resp.Header.Get("Content-Length")
	size, err := strconv.Atoi(clen)
	if err != nil {
		hlog.Fatal(err.Error())
	}

	file, err := os.Create(output)
	if err != nil {
		hlog.Fatal(err.Error())
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(
		int64(size),
		filepath.Base(url),
	)

	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		hlog.Fatal(err.Error())
	}
}
