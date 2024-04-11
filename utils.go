package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/google/uuid"
)

// A new folder is created at the root of the project.
func createFolder(dirname string) error {
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirname, 0755)
		if errDir != nil {
			return errDir
		}
	}
	return nil
}

// The mime type of the image is changed, it is compressed and then saved in the specified folder.
func imageProcessing(buffer []byte, quality int, dirname string) (string, error) {
	filename := strings.Replace(uuid.New().String(), "-", "", -1) + ".jpg"
	println(fmt.Sprintf("Compressing image %s", filename))
	image, err := vips.NewImageFromBuffer(buffer)
	if err != nil {
		return filename, err
	}
	defer image.Close()

	image.AutoRotate()

	options := vips.NewJpegExportParams()
	options.Quality = quality
    // Set MozJPEG-specific options
    options.SubsampleMode = vips.VipsForeignSubsampleOn
    options.TrellisQuant = true
    options.OvershootDeringing = true
    options.QuantTable = 3
    options.OptimizeScans = true
    options.Interlace = true
	
    imageBytes, _, _ := image.ExportJpeg(options)
	err = os.WriteFile(fmt.Sprintf("./%s/%s", dirname, filename), imageBytes, 0644)
	if err != nil {
		return filename, err
	}
	return filename, nil
}

func processDirectory(fileType, dirInput string, dirOutput string, quality int) error {
	// Ensure the output directory exists
	if err := createFolder(dirOutput); err != nil {
		return err
	}

	// Walk through the directory
	err := filepath.Walk(dirInput, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file matches the specified file type
		if !info.IsDir() && strings.HasSuffix(info.Name(), fileType) {
			// Read the file
			buffer, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Compress the image
			_, err = imageProcessing(buffer, quality, dirOutput)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

