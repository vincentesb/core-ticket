package zip_helper

import (
	"archive/zip"
	"core-ticket/base/helpers/file_helper"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ZipFiles compresses multiple files into a single zip archive file.
//
// # Requires filename parameter and list name of files to be added in zip
//
// Returns zip file with its tmp location name or errors when occur
func ZipFiles(filename string, files []string) (string, error) {
	// Ensure the directory exists
	if err := os.MkdirAll(file_helper.TmpDirectory, os.ModePerm); err != nil {
		return "", err
	}

	// Format zip tmp file
	zipFileName := fmt.Sprintf("%s/%s.zip", file_helper.TmpDirectory, filename)

	// Create a zip file
	newZipFile, err := os.Create(zipFileName)
	if err != nil {
		return "", err
	}
	defer newZipFile.Close()

	// Initialize the zip writer
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to the zip file
	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return "", err
		}
	}

	return zipFileName, nil
}

// addFileToZip adds a file to the given zip.Writer
func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	// Create a zip header based on file info
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filename) // Use file name in zip archive
	header.Method = zip.Deflate           // Set compression method

	// Create a writer for the file in the zip archive
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copy file data to zip writer
	_, err = io.Copy(writer, fileToZip)
	return err
}
