package terminal

import (
	"AutoEMU/autoEmuErrors"
	"AutoEMU/config"
	"AutoEMU/emulator"
	"AutoEMU/logger"
	"AutoEMU/manager"
	"AutoEMU/shared"

	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var err = errors.New("TMP")
var sep string = "==================================================================================================="

func emulateFromFirmware() error {

	var firmwareFile string
	var machineName string

	for {
		firmwareFile = shared.GetInputWithPrompt("Enter path to firmware file", "No path entered")
		err = getFile(firmwareFile)
		if err != nil {
			logger.Warn(err.Error(), true)
		} else {
			break
		}
	}

	for {
		machineName = shared.GetInputWithPrompt("Enter a unique machine name", "No machine name entered")
		err = getMachineName(machineName)
		if err != nil {
			logger.Warn(err.Error(), true)
		} else {
			break
		}
	}

	logger.Info("Starting Emulation Process", true)
	err := emulator.EmulateFromFirmware(firmwareFile, machineName)
	if err != nil {
		// Delete the machine
		manager.DeleteMachine(machineName)
	}
	return err
}

func importMachine() error {

	var zipFile string
	for {
		zipFile = shared.GetInputWithPrompt("Enter path to zipped machine config", "No path entered")
		err = getFile(zipFile)
		if err != nil {
			logger.Warn(err.Error(), true)
		} else if strings.Split(filepath.Base(zipFile), ".")[1] != "zip" {
			logger.Warn(
				"Not a zip file",
				true,
			)
		} else {
			break
		}
	}

	err = emulator.EmulateFromImport(zipFile)

	return err
}

func emulateFromMachine() error {
	
	// Get machine name
	var machineName string
	for {
		machineName = shared.GetInputWithPrompt("Enter machine name or -list to view configued machines", "No machine name entered")
		if machineName == "-list" {
			listMachines()
		} else {
			var found = false
			machines, _ := os.ReadDir(
				fmt.Sprintf("%s/FirmwareFiles", config.Configuration.WorkSpace),
			)
			for _, machine := range machines {
				if machineName == machine.Name() {
					found = true
					break     
				}
			}
			if found {
				logger.Info(
					fmt.Sprintf("Found machine %s", machineName),
					true,
				)
				break
			} else {
				logger.Error(
					fmt.Sprintf("Could not find machine %s", machineName),
					true,
				)
			}
		}
	}

	// Emulate machine
	err := emulator.EmulateMachine(machineName)
	
	return err
}

func getFile(file string) error {
	
	// Verify file exists
	if _, err = os.Stat(file); err == nil {
		// File found
		msg := fmt.Sprintf("File %s found\n", file)
		logger.Info(msg, false)
		return nil
	} else {
		// File not found
		return fmt.Errorf("could not find file: %s", file)
	}
}

func getMachineName(machineName string) error {
	
	// Check if name is already taken
	var machines []os.DirEntry
	machines, err = os.ReadDir(
		fmt.Sprintf("%s/FirmwareFiles", config.Configuration.WorkSpace),
	)
	for _, machine := range machines {
		if machineName == machine.Name() {
			// Machine name taken
			return fmt.Errorf("machine name %s taken", machineName)
			
		}
	}
	return nil
}

func changeConfig() {
	var field int
	fmt.Println(sep)
	fmt.Println("Select Configuration field to change")
	fmt.Println("1) WorkSpace Directory")
	fmt.Println("2) Logs Directory")
	for ok := true; ok; {
		ok = false
		input := shared.GetInput("No option selected")
		field, err = strconv.Atoi(input)
		if err != nil || field  < 1 || field > 2 {
			logger.Error("Invalid Option", true)
			ok = true
		}
	}
	config.Configuration.UpdateConfig(field)
	fmt.Println(sep)
	// Menu
	getOption()
}

func deleteMachine() error {
	var input string
	for {
		// Get Machine name
		input = shared.GetInputWithPrompt("Enter name of machine to delete", "No machine name entered")
		// Check machine exists
		pathToMachine := fmt.Sprintf(
			"%s/FirmwareFiles/%s",
			config.Configuration.WorkSpace,
			input,
		)
		logger.Info(pathToMachine, true)
		_, err := os.Stat(pathToMachine)
		if os.IsNotExist(err) {
			logger.Error(
				fmt.Sprintf("Could not find machine %s to delete", input),
				true,
			)
			return &autoEmuErrors.DeletionError{
				Machine: input,
				Err: err,
			}
		} else {
			break
		}
	}
	
	var prompt = fmt.Sprintf(
		"Are you sure you want to delete machine %s (Y/N)",
		input,
	)
	if !shared.GetInputBinary(prompt, "Y", "N") {
		logger.Info("Canceling deletion", true)
		return nil
	}
	logger.Info("Machine will now be deleted", true)
	return manager.DeleteMachine(input)
}

func listMachines() {
	fmt.Println(shared.Separator)
	fmt.Println("Configured Machines: ")
	path := fmt.Sprintf(
		"%s/FirmwareFiles",
		config.Configuration.WorkSpace,
	)
	entries, _ := os.ReadDir(path)
	for _,machine := range entries {
		fmt.Println(machine.Name())
	}
	fmt.Println(shared.Separator)
}