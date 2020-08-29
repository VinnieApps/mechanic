package tarutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndUntar(t *testing.T) {
	tempDir, tempDirErr := ioutil.TempDir("", "tartest")
	require.Nil(t, tempDirErr, "Should be able to create temp directory")

	packageDir := filepath.Join(tempDir, "package")
	os.Mkdir(packageDir, 0777)

	file1Content := "This is a test"
	ioutil.WriteFile(filepath.Join(packageDir, "test1.txt"), []byte(file1Content), 0666)

	file2Content := "This is some more content"
	file2Path := filepath.Join(packageDir, "subdir")
	os.Mkdir(file2Path, 0777)
	ioutil.WriteFile(filepath.Join(file2Path, "test2.txt"), []byte(file2Content), 0666)

	tarFile := filepath.Join(tempDir, "package.tar.gz")
	tarError := CreateTarFromDir(tarFile, packageDir)
	require.Nil(t, tarError, "Should tar correctly")

	outputDir := filepath.Join(tempDir, "output")
	files, untarErr := UntarFile(tarFile, outputDir)

	require.Nil(t, untarErr, "Should untar correctly")
	assert.Equal(t, 2, len(files), "Should untar two files")
	assert.Contains(t, files, "test1.txt", "Should match tarred file names")
	assert.Contains(t, files, "subdir/test2.txt", "Should match tarred file names")
}
