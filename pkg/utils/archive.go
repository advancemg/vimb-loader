package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func RecursiveZip(outputZip, inputPath string) error {
	destinationFile, err := os.Create(inputPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(outputZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(outputZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}
