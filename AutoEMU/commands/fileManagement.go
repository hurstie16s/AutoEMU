package commands

import (
	"AutoEMU/config"
	"AutoEMU/logger"
	"fmt"
	"os/exec"
)

func CreateDIRFromWorkSpace(path string) error {
	cmd := exec.Command(
		"mkdir",
		path,
	)
	cmd.Dir = config.Configuration.WorkSpace
	msg := fmt.Sprintf("Dir set to: %s", config.Configuration.WorkSpace)
	logger.Info(msg, true)
	err := cmd.Run()
	if err != nil {
		msg := fmt.Sprintf("Failed to create machine directory: %v", err.Error())
		logger.Error(msg, true)
	}
	return err
}

func CopyFile(from string, to string) error {

	cmd := exec.Command(
		"cp",
		from,
		to,
	)
	cmd.Dir = config.Configuration.WorkSpace
	_, err := cmd.Output()
	return err
}

func MoveFile(source string, dest string) bool {
	cmd := exec.Command(
		"mv",
		source,
		dest,
	)
	cmd.Dir = config.Configuration.WorkSpace
	_, err := cmd.Output()
	return err == nil
}

func ReadLink() (error, string) {
	cmd := exec.Command("readlink", "-f", "FirmwareFiles/")
	cmd.Dir = config.Configuration.WorkSpace
	out, err := cmd.Output()
	return err, string(out)
}
