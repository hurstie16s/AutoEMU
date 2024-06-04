package terminal

import (
	"AutoEMU/logger"
	"AutoEMU/shared"
	
	"fmt"
	"os"
	"strings"
)

var title = fmt.Sprint(
	"    _         _        _____ __  __ _   _ \n" +
		"   / \\  _   _| |_ ___ | ____|  \\/  | | | |\n" +
		"  / _ \\| | | | __/ _ \\|  _| | |\\/| | | | |\n" +
		" / ___ \\ |_| | || (_) | |___| |  | | |_| |\n" +
		"/_/   \\_\\__,_|\\__\\___/|_____|_|  |_|\\___/ ",
)
var titleOut = fmt.Sprintf("\n%s\n%sAutomated Emulation Platform%s\n", title, titlePadding, titlePadding)
var menu = fmt.Sprint(
	"1) Emulate from Firmware\n" +
		"2) Import Configuration\n" +
		"3) Run configured machine\n" +
		"4) Change Configuration\n" +
		"5) Delete Machine\n" +
		"6) List configured machines\n" +
		"Enter choice, Q to quit of H for help\n",
)

const titlePadding = "       "

func LaunchBlank() {
	getOption()
}

func getOption() {

	fmt.Print(titleOut)
	fmt.Print(menu)

	var option = shared.GetInput("No option entered")
	msg := fmt.Sprintf("Option Selected: %v", option)

	logger.Info(msg, false)
	handleOption(option)
}

func handleOption(option string) {

	switch optionParse := strings.ToUpper(option); optionParse {
	case "Q":
		logger.Info("Terminating Application", true)
		os.Exit(1)
	case "H":
		logger.Info("Printing help page", false)
		printFullHelp()
		printVSpace(3)
		getOption()
	case "1":
		logger.Info("Emulating device from firmware", true)
		// Prep to emulate a device from firmware file
		err = emulateFromFirmware()
		if err != nil {
			var msg = fmt.Sprintf("Fatal Error in firmware emulation: %s", err.Error())
			logger.Error(msg, true)
			printVSpace(1)
			getOption()
		}
	case "2":
		logger.Info("Importing machine configuration zip file", true)
		err = importMachine()
		if err != nil {
			logger.Error(
				fmt.Sprintf("Fatal Error importing machine: %v", err),
				true,
			)
			printVSpace(1)
			getOption()
		}
	case "3":
		logger.Info("Emulating from built machine", true)
		err = emulateFromMachine()
		if err != nil {
			logger.Error(
				fmt.Sprintf("Fatal error in emulation process: %v", err),
				true,
			)
			printVSpace(1)
			getOption()
		}
	case "4":
		logger.Info("Changing configuration", false)
		changeConfig()
	case "5":
		logger.Info("Deleting machine", true)
		err := deleteMachine()
		if err != nil {
			logger.Error(err.Error(), true)
		}
	case "6":
		listMachines()
	default:
		msg := fmt.Sprintf("%v is not a valid option\n", option)
		logger.Warn(msg, true)
		getOption()
	}
	getOption()
}
