package main

// System Imports
import (
	"AutoEMU/config"
	"AutoEMU/logger"
	"AutoEMU/terminal"
	"fmt"
	"log"
	"os"
	"time"
)

//"iot-virtualiser/emulator"

/*
- architercture
- machine/ board
- cpu
- memory
- disk
- rfs
- kernal
- os
- monitor options
- full disk image/ snapshot
*/

var err error

func main() {

	// Get Command Line Arguments
	argsWithProg := os.Args
	args := argsWithProg[1:]

	var configFlag = len(args) != 0 && args[0] == "-config"

	// Load Configuration
	err = config.LoadConfig(configFlag)

	// Non interactive config
	if configFlag {
		args = terminal.ProcessConfig(args)
	}

	if err != nil {
		msg := fmt.Sprintf("Configuration load failed: %v", err.Error())
		fmt.Printf("%s%s", logger.FatalMessage, msg)
		os.Exit(1)
	}

	config.CheckWorkSpace()

	// Setup logger
	// Create log file and set output
	var logFile *os.File
	var logFileName = fmt.Sprintf("%s/log-%d.log", config.Configuration.LogDir, time.Now().Nanosecond())
	logFile, err = os.Create(logFileName)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	logger.Config("Log Output Set", false)
	
	if len(args) == 0 && !configFlag {
		terminal.LaunchBlank()
	} else if len(args) != 0 {
		terminal.ProcessArgs(args)
	} else {
		os.Exit(0)
	}

	// sudo apt install bc binutils bison dwarves flex gcc git gnupg2 gzip libelf-dev libncurses5-dev libssl-dev make openssl pahole perl-base rsync tar xz-utils

	os.Exit(0)

}
