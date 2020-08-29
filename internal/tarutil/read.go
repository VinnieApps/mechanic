package tarutil

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// UntarFile will decompress and copy the files from a gzipped tar file into
// the output directory. A list of the path for each file decompressed and
// copied will be returned.
func UntarFile(tarFilePath string, outputDir string) ([]string, error) {
	file, openErr := os.Open(tarFilePath)
	if openErr != nil {
		return nil, openErr
	}
	defer file.Close()

	gzf, unZipErr := gzip.NewReader(file)
	if unZipErr != nil {
		return nil, unZipErr
	}

	tarReader := tar.NewReader(gzf)
	files := make([]string, 0)
	for {
		header, readHeaderErr := tarReader.Next()
		if readHeaderErr == io.EOF {
			break
		}

		if header.Typeflag == tar.TypeReg {
			writeTo, absErr := filepath.Abs(filepath.Join(outputDir, header.Name))
			if absErr != nil {
				return nil, absErr
			}

			dirToCreate := filepath.Dir(writeTo)
			if err := os.MkdirAll(dirToCreate, 0755); err != nil {
				return nil, err
			}

			files = append(files, header.Name)

			fileToWrite, fileOpenErr := os.OpenFile(writeTo, os.O_RDWR|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode())
			if fileOpenErr != nil {
				return nil, fileOpenErr
			}

			if _, copyErr := io.Copy(fileToWrite, tarReader); copyErr != nil {
				return nil, copyErr
			}

			if err := fileToWrite.Close(); err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}
