package manager

import (
	"AutoEMU/autoEmuErrors"
	"AutoEMU/config"
	"AutoEMU/logger"
	"AutoEMU/machine"

	"archive/zip"
	"strconv"
	"errors"
	"fmt"
	"os"
	"io"
)

func ExportMachine(machine machine.Machine) error {

	// zipFilePath
	path := fmt.Sprintf(
		"%s/ZippedMachines/%s.zip",
		config.Configuration.WorkSpace,
		machine.Name,
	)
	// ensure file is unique
	path = ensureUnique(path)
	// Create zip file
	zipFile, err := os.Create(path)
	if err != nil {
		return &autoEmuErrors.MachineExportError{
			Machine: machine.Name,
			Err: fmt.Errorf("failed to create zip file"),
		}
	}
	defer zipFile.Close()

	// Create writer for zip file
	zipper := zip.NewWriter(zipFile)
	defer zipper.Close()

	// Make array of files that will be zipped
	/*
	Firmware File
	ISO Image
	Hard Drive Image
	Machine Config yaml File
	*/
	prefix := fmt.Sprintf(
		"%s/FirmwareFiles/%s/",
		config.Configuration.WorkSpace,
		machine.Name,
	)
	var files = [4]string{
		machine.MachineFiles.FirmwareFile,
		machine.MachineFiles.Kernel,
		machine.MachineFiles.HardDriveFile,
		fmt.Sprintf("%s.yaml", machine.Name),

	}
	// Loop through files, writing each file to zipper
	for _, file := range files {
		// Open file
		fileForZipper, err := os.Open(
			fmt.Sprintf(
				"%s%s",
				prefix,
				file,
			),
		)
		if err != nil {
			return &autoEmuErrors.MachineExportError{
				Machine: machine.Name,
				Err: fmt.Errorf("failed to open file %s for zipping", file),
			}
		}
		defer fileForZipper.Close()

		// Make new header for file
		info, err := fileForZipper.Stat()
		if err != nil {
			return &autoEmuErrors.MachineExportError{
				Machine: machine.Name,
				Err: fmt.Errorf("failed to get information for file %s", file),
			}
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return &autoEmuErrors.MachineExportError{
				Machine: machine.Name,
				Err: fmt.Errorf("failed to make new header for file %s", file),
			}
		}

		// Set header name
		header.Name = file

		// Add file header to archive
		writer, err := zipper.CreateHeader(header)
		if err != nil {
			return &autoEmuErrors.MachineExportError{
				Machine: machine.Name,
				Err: fmt.Errorf("failed to add header to archive for file %s", file),
			}
		}

		// Write file contents to the archive.
        _, err = io.Copy(writer, fileForZipper)
        if err != nil {
            return &autoEmuErrors.MachineExportError{
				Machine: machine.Name,
				Err: fmt.Errorf("failed to write file %s to archive", file),
			}
        }
	}
	logger.Success(
		fmt.Sprintf("Machine %s exported to zip file", machine.Name),
		true,
	)
	return nil
}

func ensureUnique(path string) string {
	var counter int = 2
	var ammendAttempted bool = false
	for {
		if _, err := os.Stat(path); err == nil {
			logger.Warn(
				fmt.Sprintf("File %s already exists, amending file", path),
				true,
			)
			if ammendAttempted {
				// Remove last 5 chars of file path (num.zip)
				path = path[:len(path)-5]
				// Create new path
				path = fmt.Sprintf(
					"%s%s.zip",
					path,
					strconv.Itoa(counter),
				)
				counter += 1
			} else {
				// Remove last 4 chars of file path (num.zip)
				path = path[:len(path)-4]
				// Create new path
				path = fmt.Sprintf(
					"%sv%s.zip",
					path,
					strconv.Itoa(counter),
				)
				counter += 1
			}
		} else if errors.Is(err, os.ErrNotExist) {
			return path
		} else {
			logger.Warn("Cannot determine if file already exists, file may be overridden", true)
			return path
		}
	}
}