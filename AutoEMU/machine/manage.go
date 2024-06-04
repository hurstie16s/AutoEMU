package machine

import (
	"AutoEMU/autoEmuErrors"
	"AutoEMU/config"
	"AutoEMU/logger"

	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func (machine *Machine) Export() error {

	var machineFile *os.File
	var err error

	// Create file path
	filePath := fmt.Sprintf(
		"%s/FirmwareFiles/%s/%s.yaml", 
		config.Configuration.WorkSpace, 
		machine.Name, machine.Name,
	)

	msg := fmt.Sprintf("Exporting machine %s to %s", machine.Name, filePath)
	logger.Info(msg, true)
	logger.Warn("If this file exists, it will be overwritten", true)

	// Open/ create file
	machineFile, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return &autoEmuErrors.MachineExportError{
			Machine: machine.Name, 
			Err: fmt.Errorf("failed to open file %s: %s", filePath, err),
		}
	}
	defer machineFile.Close()
	
	// Create encoder
	encoder := yaml.NewEncoder(machineFile)
	defer encoder.Close()

	// Write machine to file
	err = encoder.Encode(machine)
	if err != nil {
		return &autoEmuErrors.MachineExportError{
			Machine: machine.Name,
			Err: fmt.Errorf("failed to write machine configuration to %s: %v", machineFile.Name(), err),
		}
	}
	
	return nil
}