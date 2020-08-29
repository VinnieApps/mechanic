package fileutil

import (
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file from source to destination
func CopyFile(srcPath string, destPath string) error {
	parent := filepath.Dir(destPath)
	if err := os.MkdirAll(parent, 0700); err != nil {
		return err
	}

	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
