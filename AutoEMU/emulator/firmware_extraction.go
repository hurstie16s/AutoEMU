package emulator

// Imports
import (
	"AutoEMU/autoEmuErrors"
	"AutoEMU/commands"
	"AutoEMU/config"
	"AutoEMU/logger"
	"AutoEMU/machine"
	"io"

	"errors"
	"fmt"
	"os/exec"
	"strings"
)

//var extractorScript string = fmt.Sprintf("%s/FirmwareMod/extract_firmware.sh", config.Configuration.WorkSpace)

/*
Given a file or zip file and a name for the firmware
if zip file, unzip it
else, just rename
into its own dir name_foler
run ./extract_firmware.sh on the firmware file
convert rfs img file into .qcow2 file
*/

func ExtractFirmware(pathToFile string, name string) (machine.Machine, error) {

	var err error
	var machineConfig machine.Machine

	// Move firmware into workspace
	machineConfig, err = moveFileToWorkSpace(pathToFile, name)
	if err != nil {
		return machineConfig, &autoEmuErrors.FileRelocationError{Err: err}
	}

	// Run ./extract_firmware.sh on file
	target := fmt.Sprintf("FirmwareFiles/%s/%s", machineConfig.Name, machineConfig.MachineFiles.FirmwareFile)
	logger.Info(
		fmt.Sprintf("Extraction target: %s", target),
		true,
	)

	err, _ = commands.ReadLink()
	if err != nil {
		return machineConfig, &autoEmuErrors.ReadLinkError{Err: err}
	}

	var permsChecked = false
	for {
		logger.Info("Running Extraction Script", true)

		cmd := exec.Command("./FirmwareMod/extract-firmware.sh", target) // TODO: Alter outputs to suit what I need

		cmd.Dir = config.Configuration.WorkSpace

		var out []byte
		out, err = cmd.Output()
		logger.Write(string(out), false) // TODO: get rfs and add to machineConfig
		if err == nil {
			logger.Success("Firmware Extracted", true)
			// TODO: Check firmware has been extracted
			break
		} else {
			switch err.Error() {
			case "exit status 2":
				logger.Error("fmk directory found, attempting to delete", true)
				cmd := exec.Command(
					"rm",
					"-r",
					"fmk",
				)
				cmd.Dir = config.Configuration.WorkSpace
				err = cmd.Run()
				if err != nil {
					logger.Error(
						fmt.Sprintf("Failed to delete fmk directory: %v", err),
						true,
					)
					return machineConfig, &autoEmuErrors.FirmwareExtractionError{Err: err}
				}
				logger.Success("fmk directory deleted, re-running extraction script", true)
			case "exit status 3":
				return machineConfig, &autoEmuErrors.FirmwareExtractionError{
					Err: fmt.Errorf("no supported file system could be detected for machine %s", machineConfig.Name),
				}
			default:
				logger.Error("Failed to extract firmware", true)
				if permsChecked {
					logger.Warn("File permissions already checked", true)
					return machineConfig, &autoEmuErrors.FirmwareExtractionError{Err: err}
				} else {
					checkPerms("FirmwareMod")
					permsChecked = true
					logger.Info("Re-running extraction script", true)
				}	
			}	
		}		
	}
	// Move fmk file into the firmware location
	var source = "fmk"
	dest := fmt.Sprintf("FirmwareFiles/%s", machineConfig.Name)
	commands.MoveFile(source, dest)
	logger.Info("Extracted fmk directory moved into machine directory", false)

	/*
	var imageFile = "fmk/image_parts/rootfs.img" // In future this may need to be variable

	msg := fmt.Sprintf("Building hard drive image %s.qcow2", machineConfig.Name)
	logger.Info(msg, true)

	hardDriveImage, err = ConvertImage(imageFile, machineConfig.Name, "qcow2")
	if err != nil {
		return machineConfig, &autoEmuErrors.RFSConversionError{Err: err}
	}
	machineConfig.MachineFiles.HardDriveFile = hardDriveImage

	logger.Success("Hard drive image created", true)
	*/
	return machineConfig, nil

}

func moveFileToWorkSpace(pathToFile string, name string) (machine.Machine, error) {

	var machineConfig machine.Machine = machine.Machine{}
	machineConfig.Name = name

	// Get File extension
	index := strings.LastIndex(pathToFile, ".")
	extension := string(pathToFile[len(pathToFile)-(len(pathToFile)-index)+1:])
	machineConfig.FirmwareExtension = extension
	msg := fmt.Sprintf("Extension Found: %s\n", extension)
	logger.Info(msg, false)

	// Create DIR for firmware
	dir := fmt.Sprintf("FirmwareFiles/%s", name)
	err := commands.CreateDIRFromWorkSpace(dir)
	if err != nil {
		return machineConfig, fmt.Errorf("failed to create directory in workspace: %v", err)
	}
	logger.Info("Directory created in workspace", true)

	// Append extension to filename
	file := fmt.Sprintf("%s.%s", name, extension)
	machineConfig.MachineFiles.FirmwareFile = file

	// Copy file to new location
	newFile := fmt.Sprintf("%s/%s", dir, file)
	err = commands.CopyFile(pathToFile, newFile)
	if err != nil {
		return machineConfig, errors.New("failed to move firmware file into workspace")
	}
	return machineConfig, nil
}

func checkPerms(dirInWorkSpace string) {

	logger.Info(
		fmt.Sprintf("Checking file permissions for %s", dirInWorkSpace),
		true,
	)

	// Ensure the actual dir is set to 755
	cmd := exec.Command(
		"chmod",
		"755",
		dirInWorkSpace,
	)
	cmd.Dir = config.Configuration.WorkSpace
	err = cmd.Run()
	if err != nil {
		logger.Error(err.Error(), true)
		return
	}

	err = fixPerms("d", fmt.Sprintf("%s/%s", config.Configuration.WorkSpace, dirInWorkSpace))
	if err != nil {
		logger.Error("Failed to fix file permissions", true)
		return
	}

	err = fixPerms("f", fmt.Sprintf("%s/%s", config.Configuration.WorkSpace, dirInWorkSpace))
	if err != nil {
		logger.Error(err, true)
		return
	}
}

func fixPerms(typeForPerms string, dir string) error {
	
	cmd1 := exec.Command(
		"find",
		".",
		"-type",
		typeForPerms,
		"-print0",
	)
	cmd1.Dir = dir

	cmd2 := exec.Command(
		"xargs",
		"-0",
		"chmod",
		"755",
	)
	cmd2.Dir = dir

	read, write := io.Pipe()
	cmd1.Stdout = write
	cmd2.Stdin = read

	cmd1.Start()
	cmd2.Start()
	cmd1.Wait()
	write.Close()
	err := cmd2.Wait()

	return err
}

//func extractScript
