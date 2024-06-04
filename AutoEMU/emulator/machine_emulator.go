package emulator

import (
	"AutoEMU/autoEmuErrors"
	"AutoEMU/commands"
	"AutoEMU/config"
	"AutoEMU/logger"
	"AutoEMU/machine"
	"AutoEMU/manager"
	"AutoEMU/shared"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func EmulateFromFirmware(pathToFirmware string, name string) error {
	var machineConfig machine.Machine
	var err error

	logger.Info("Starting Firmware Extraction", true)

	machineConfig, err = ExtractFirmware(pathToFirmware, name)
	if err != nil {
		return err
	}

	machineConfig = analyseConfigLog(machineConfig)

	// Try to automatically get architecture information
	machineConfig = GetArchInfo(machineConfig)

	/*
		// Build ISO image - may be unessecary
		machineConfig, err = MakeISO(machineConfig)
		if err != nil {
			return err
		}
	*/

	// Analyse what data was collected from architecture analysis
	machineConfig = checkAnalysis(machineConfig)

	// Get kernel & rfs
	machineConfig = getKernel(machineConfig)

	machineConfig, _, flag := System(machineConfig)

	// Write Machine configuration to yaml file
	machineConfig.Export()

	// Check if user wishes to export machine
	
	exportMachine := shared.GetInputBinary(
		"Do you wish to export the machine to a zip file (Y/N)",
		"Y",
		"N",
	)
	// done for overnight tesitng
	if exportMachine {
		manager.ExportMachine(machineConfig)
	}

	if flag {
		// Run script
		logger.Info("Running Script", true)
		runScript(
			fmt.Sprintf("%s/FirmwareFiles/%s", config.Configuration.WorkSpace, machineConfig.Name),
			machineConfig.MachineFiles.ShellScript,
		)
	}
	fmt.Scanln()

	return nil

	// search root file system
	/*
		Take path to firmware file
		Determine Arch, machine and other nessecary extractions
		Extract firmware and build config
		Use generated config to buid qemu command
		Run command - Within seperate go routine
		Use libvirt and qemu monitor to communicate with emulated machine
		Snap qemu window into UI (if can)
		Save config to a yaml file in /WorkSpace/Machines
		Give user name of config
	*/
}

func EmulateFromImport(pathToZip string) error {
	// Run shell script
	// Read creds from yaml file

	// Get file name
	file := filepath.Base(pathToZip)
	fileName := strings.Split(file, ".")[0]
	// Create dir to extract zip file into
	commands.CreateDIRFromWorkSpace(
		fmt.Sprintf("FirmwareFiles/%s", fileName),
	)
	// Copy zip file into
	commands.CopyFile(
		pathToZip,
		fmt.Sprintf("FirmwareFiles/%s/%s", fileName, file),
	)

	// Unzip the file
	cmd := exec.Command(
		"unzip",
		file,
	)
	cmd.Dir = fmt.Sprintf(
		"%s/FirmwareFiles/%s",
		config.Configuration.WorkSpace,
		fileName,
	)
	err = cmd.Run()
	if err != nil {
		return &autoEmuErrors.UnzipError{
			FileName: file,
			Err:      err,
		}
	}

	// Remove zip file
	cmd = exec.Command(
		"rm",
		"*.zip",
	)
	cmd.Dir = fmt.Sprintf(
		"%s/FirmwareFiles/%s",
		config.Configuration.WorkSpace,
		fileName,
	)
	cmd.Run()
	logger.Success("Machine config unzipped", true)

	logger.Success("Machine Imported", true)

	return nil
}

func EmulateMachine(machineName string) error {

	machine, err := machine.Load(machineName)
	if err != nil {
		return err
	}
	logger.Success(
		fmt.Sprintf("Loaded machine %s into AutoEMU", machine.Name),
		true,
	)

	// Check if user wishes to export machine

	exportMachine := shared.GetInputBinary(
		"Do you wish to export the machine to a zip file (Y/N)",
		"Y",
		"N",
	)

	// done for overnight tesitng
	if exportMachine {
		manager.ExportMachine(machine)
	}

	if machine.MachineFiles.ShellScript != "" {
		logger.Info("Running shell script containing qemu command", true)
		logger.Warn("This script may not have been built by AutoEMU, the full contents of the script will be run", true)
		fmt.Printf("Enter 'RUN' to run script %s, and anything else to abort\n>>>", machine.MachineFiles.ShellScript)
		var check string
		fmt.Scanln(&check)
		if check != "RUN" {
			logger.Info("Aborting emulation", true)
			return nil
		}
		logger.Info("Running Script", true)

		runScript(
			fmt.Sprintf("%s/FirmwareFiles/%s", config.Configuration.WorkSpace, machine.Name),
			machine.MachineFiles.ShellScript,
		)

	} else {
		machine, _, flag := System(machine)
		if flag {
			// Run shell script
			logger.Info("Running Shell Script", true)
			runScript(
				fmt.Sprintf("%s/FirmwareFiles/%s", config.Configuration.WorkSpace, machine.Name),
				machine.MachineFiles.ShellScript,
			)
		}
	}

	return nil
}
