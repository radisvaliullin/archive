package dbdump

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FolderToZIP - archives folder files to zip.
func FolderToZIP(folderpath, zippath string, del bool) error {

	// create zip file
	zipDir := filepath.Dir(zippath)
	if _, err := os.Stat(zipDir); os.IsNotExist(err) {
		if err := os.MkdirAll(zipDir, os.ModePerm); err != nil {
			return err
		}
	}
	zipfile, err := os.Create(zippath)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	// zip file writer
	zipw := zip.NewWriter(zipfile)
	defer zipw.Close()

	// zipped directory
	zipedDir := filepath.Base(folderpath)

	// zip folder files add to zip
	filepath.Walk(folderpath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileHeader, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}
		fileHeader.Name = filepath.Join(zipedDir, strings.TrimPrefix(path, folderpath))

		if fileInfo.IsDir() {
			fileHeader.Name += "/"
		} else {
			fileHeader.Method = zip.Deflate
		}

		filew, err := zipw.CreateHeader(fileHeader)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		// Archived file
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(filew, f)
		return err
	})

	if del {
		err := os.RemoveAll(folderpath)
		if err != nil {
			return err
		}
	}

	return err
}
