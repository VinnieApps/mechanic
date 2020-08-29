package work

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/vinnieapps/mechanic/internal/fileutil"
)

// EnsureFileTask will make sure the target file exists and
// has the same contents as the source file
type EnsureFileTask struct {
	Source string
	Target string

	packageBase string
	sourceFile  string
}

// Execute executes the ensure file task
func (t *EnsureFileTask) Execute(packageBase string) error {
	t.packageBase = packageBase
	t.sourceFile = filepath.Join(packageBase, t.Source)

	if err := t.validate(); err != nil {
		return err
	}

	sourceContent, readSourceErr := ioutil.ReadFile(t.sourceFile)
	if readSourceErr != nil {
		return readSourceErr
	}

	targetContent, readTargetErr := t.getTargetFileContent()
	if !os.IsNotExist(readTargetErr) {
		return readTargetErr
	}

	if targetContent != nil {
		sourceSha := base64.StdEncoding.EncodeToString(sha256.New().Sum(sourceContent))
		targetSha := base64.StdEncoding.EncodeToString(sha256.New().Sum(targetContent))
		if sourceSha == targetSha {
			fmt.Printf("Source file '%s' has the same content as target '%s'\n", t.sourceFile, t.Target)
			return nil
		}
	}

	fmt.Printf("Copying source file '%s' to target '%s'\n", t.sourceFile, t.Target)
	return fileutil.CopyFile(t.sourceFile, t.Target)
}

func (t *EnsureFileTask) getTargetFileContent() ([]byte, error) {
	stat, statErr := os.Stat(t.Target)
	if statErr != nil {
		return nil, statErr
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("'%s' is a directory", t.Target)
	}

	return ioutil.ReadFile(t.Target)
}

func (t *EnsureFileTask) validate() error {
	if t.sourceFile == "" || t.Target == "" {
		return fmt.Errorf("Ensure file task needs source and target files. Target: '%s', Source: '%s", t.Target, t.sourceFile)
	}

	if _, err := os.Stat(t.sourceFile); os.IsNotExist(err) {
		return fmt.Errorf("Source file does not exist: %s", t.sourceFile)
	}

	return nil
}
