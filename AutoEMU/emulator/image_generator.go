package emulator

// Imports
import (
	"AutoEMU/config"
	"AutoEMU/logger"

	"fmt"
	"os/exec"
)

func GenImage(format string, name string, size int16) bool {

	// Command and Tags
	baseCommand := "qemu-img"
	task := "create"
	tag := "-f"
	imageType := format
	fileName := fmt.Sprintf("%s.%s", name, format)
	fileSize := fmt.Sprintf("%d.G", size)

	cmd := exec.Command(
		baseCommand,
		task,
		tag,
		imageType,
		fileName,
		fileSize,
	)
	// Set Execution Directory
	cmd.Dir = fmt.Sprintf("%s/FirmwareFiles/%s", config.Configuration.WorkSpace, name)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true

}

func ConvertImage(imageFile string, machineName string, format string) (string, error) {

	var outputFile string = fmt.Sprintf("%s.%s", machineName, format)

	cmd := exec.Command(
		"qemu-img",
		"convert",
		"-O",
		format,
		imageFile,
		outputFile,
	)
	cmd.Dir = fmt.Sprintf("%s/FirmwareFiles/%s", config.Configuration.WorkSpace, machineName)
	err := cmd.Run()
	if err != nil {
		logger.Error("Handle", true)
	}
	return fmt.Sprintf("%s.%s", machineName, format), nil
}

/*
func MakeISO(machine machine.Machine) (machine.Machine, error) {

	workingDir := fmt.Sprintf("%s/FirmwareFiles/%s", config.Configuration.WorkSpace, machine.Name)

	cmd := exec.Command(
		"genisoimage",
		"-o",
		fmt.Sprintf("%s.iso", machine.Name),
		"fmk/rootfs",
	)
	cmd.Dir = workingDir

	out, err := cmd.Output()
	if err != nil {
		return machine, &autoEmuErrors.ISOGenerationError{Err: err}
	}
	logger.Write(string(out), false)
	msg := fmt.Sprintf("ISO image %s.iso created", machine.Name)
	logger.Success(msg, true)

	machine.MachineFiles.ISOImage = fmt.Sprintf("%s.iso", machine.Name)

	return machine, nil
}
*/