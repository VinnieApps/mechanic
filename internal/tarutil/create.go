package tarutil

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// CreateTarFromDir creates a tar file including all files from a directory
func CreateTarFromDir(outputFile string, basePath string) error {
	tarFile, createFileErr := os.Create(outputFile)
	if createFileErr != nil {
		return createFileErr
	}
	defer tarFile.Close()

	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return addFileToTarWriter(path, basePath, tarWriter)
	})
}

func addFileToTarWriter(filePath string, relativeTo string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error while openning file '%s': '%s'", filePath, err.Error())
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Error while getting file info '%s': '%s'", filePath, err.Error())
	}

	relativePath, relErr := filepath.Rel(relativeTo, filePath)
	if relErr != nil {
		return relErr
	}

	log.Printf("Adding %s as %s\n", filePath, relativePath)

	header := &tar.Header{
		Name:    relativePath,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("Error while writing tar header '%s': '%s'", filePath, err.Error())
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return fmt.Errorf("Error while writing '%s' to the tarball: '%s'", filePath, err.Error())
	}

	return nil
}
