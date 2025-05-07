package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func installPackage(src, dst string, extractFn func(string, string) error) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	if extractFn != nil {
		if err := extractFn(src, dst); err != nil {
			return fmt.Errorf("extract %s: %w", src, err)
		}
	} else {
		base := filepath.Base(src)
		dstPath := filepath.Join(dst, base)
		if err := copyFile(src, dstPath); err != nil {
			return fmt.Errorf("copy %s: %w", src, err)
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	if _, err := io.Copy(df, sf); err != nil {
		return err
	}
	if fi, err := sf.Stat(); err == nil {
		os.Chmod(dst, fi.Mode())
	}
	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		tgt := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(tgt, info.Mode())
		}
		return copyFile(path, tgt)
	})
}
