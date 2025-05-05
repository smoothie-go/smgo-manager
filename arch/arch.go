package arch

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/smoothie-go/smgo-manager/hlog"
)

func UntarGz(source, target string) error {
	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(target, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		default:
			hlog.Error(fmt.Sprintf("Skipping: %s (unsupported type)\n", header.Name))
		}
	}

	return nil
}

func Unzip(source, target string) error {
	if s, err := os.Stat(target); err != nil || s.IsDir() != true {
		if err != nil {
			return err
		}
		return errors.New("Unable to unzip " + source + " as target isn't a directory")
	}

	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	hlog.Ok("Extracting " + source + "into" + target)
	defer reader.Close()

	for _, f := range reader.File {
		file, err := f.Open()
		if err != nil {
			return err
		}

		filePath := filepath.Join(target, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, f.Mode())
		} else {
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return err
			}
			fc, err := os.Create(filePath)
			if err != nil {
				return err
			}

			if _, err := io.Copy(fc, file); err != nil {
				return err
			}
		}

	}
	return nil
}
