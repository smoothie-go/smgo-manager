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

	sz "github.com/bodgit/sevenzip"
	"github.com/smoothie-go/smgo-manager/hlog"
	"github.com/ulikunitz/xz"
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
		case tar.TypeSymlink:
			linkPath := filepath.Join(target, header.Name)
			target := header.Linkname
			if err := os.Symlink(target, linkPath); err != nil {
				hlog.Error(fmt.Sprintf("Warning: failed to create symlink %s -> %s: %v", linkPath, target, err))
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

			if err := os.Chmod(targetPath, header.FileInfo().Mode()); err != nil {
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

// `UntarXz(string, string) error` was vibe coded
func UntarXz(source, target string) error {
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer f.Close()

	xzr, err := xz.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to create xz reader: %w", err)
	}

	tr := tar.NewReader(xzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar header: %w", err)
		}

		targetPath := filepath.Join(target, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directories for %s: %w", targetPath, err)
			}
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write to file %s: %w", targetPath, err)
			}
			outFile.Close()
		default:
			hlog.Error(fmt.Sprintf("Skipping: %s (unsupported type)", header.Name))
		}
	}

	return nil
}

func Un7z(source, target string) {
	extractFile := func(file *sz.File) {
		filePath := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			return
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			hlog.Fatal(err.Error())
		}

		outFile, err := os.Create(filePath)
		if err != nil {
			hlog.Fatal(err.Error())
		}
		defer outFile.Close()

		f, err := file.Open()
		if err != nil {
			hlog.Fatal(err.Error())
		}
		defer f.Close()

		_, err = io.Copy(outFile, f)
		if err != nil {
			hlog.Fatal(err.Error())
		}
	}

	func() {
		archive, err := sz.OpenReader(source)
		if err != nil {
			hlog.Fatal(err.Error())
		}
		defer archive.Close()

		for _, f := range archive.File {
			extractFile(f)
		}
	}()
}
