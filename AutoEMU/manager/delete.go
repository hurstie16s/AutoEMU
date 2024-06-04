package manager

import (
	"AutoEMU/autoEmuErrors"
	"AutoEMU/config"
	"os/exec"

	"fmt"
)

func DeleteMachine(machineName string) error {

	cmd := exec.Command(
		"rm",
		"-r",
		machineName,
	)
	cmd.Dir = fmt.Sprintf(
		"%s/FirmwareFiles",
		config.Configuration.WorkSpace,
	)
	err := cmd.Run()

	if err != nil {
		return &autoEmuErrors.DeletionError{
			Machine: machineName,
			Err: err,
		}
	}
	return nil
}