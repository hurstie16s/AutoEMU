package terminal

import (
	"fmt"
)

/*
Tags
-h [tag]: show help page, full or for specific tag
-e [filename] [machineName] emulate from firmware binary, takes 2 arguments, a path to firmware file and name for firmware
-full: sets the emulation type to full system (This is the default)
-user: sets the emulation type to user for emulating linux binaries
-o [filename]: sets a file to export the config from
-i [filename] imports a configuration
-m [machineName] loads a build machine for emulation
*/

const helpStart = "AutoEMU Help Page\n"
const tagNoTagDif = "AutoEMU can be run with or without any tags.\n" +
	"Running without tags will take you to the interactive terminal application\n"
const hTag = "-h\n" +
	" Print help page\n" +
	" Takes 1 optional argument to only display the help page for a specific tag\n" +
	" If no argument is present, the full help page will be printed\n" +
	" Argument options:\n" +
	"  h: print help page for -h tag\n" +
	"  e: show help page for -e tag\n" +
	"  full: show help page for -full tag\n" +
	"  user: show help page for -user tag\n" +
	"  o: show help page for -o tag\n" +
	"  config: show help page for -config tag\n"
const eTag = "-e\n" +
	" Emulating a device from a firmware binary\n" +
	" Takes 2 Arguments, filename and machineName\n" +
	"  filename:\n" +
	"   Path to firmware binary\n" +
	"  machineName:\n" +
	"   Name for machine being emulated. This allows a machine to be reloaded in the future\n"
const fullTag = "-full\n" +
	" Emulating a full device with specific architecture.\n" +
	" Mostly firmware can be automatically analysed by the tool to determine the required architectures.\n" +
	" On occasions, user will have to specify the architecture.\n"
const userTag = "-user\n" +
	" TODO: Fill in\n"
const oTag = "-o\n" +
	" Set an output file to save a machine configuration to.\n" +
	" Takes 1 argument, path to output file\n" +
	" If the file does not exist, it will be created.\n" +
	" If the file exists, it will be overwritten.\n"
const configTag = "-config\n" +
	" Rebuild the configuration for AutoEMU\n" +
	" Takes 2 arguments\n" +
	"  Path to the workspace directory (including /Workspace) or -pwd if the /WorkSpace directory is contained within the current working directory\n" +
	"  Path to create a Logs directory or -pwd if you wish for the Logs directory to be created in the current working directory\n" +
	" WARNING: This wil erase the previous configuration\n"

func printFullHelp() {
	// Print Intro
	printVSpace(1)
	fmt.Print(helpStart)
	printVSpace(1)
	// Tag brief description
	fmt.Print(tagNoTagDif)
	printVSpace(1)
	fmt.Println("Tags:")
	fmt.Print(hTag)
	//fmt.Print(eTag)
	//fmt.Print(fullTag)
	//fmt.Print(userTag)
	//fmt.Print(oTag)
	fmt.Print(configTag)
}

func printVSpace(r int) {
	for i := 0; i < r; i++ {
		fmt.Print("\n")
	}
}

func printByTag(tag string) {
	switch tag {
	case "h":
		fmt.Print(hTag)
	case "e":
		fmt.Print(eTag)
	case "full":
		fmt.Print(fullTag)
	case "user":
		fmt.Print(userTag)
	case "o":
		fmt.Print(oTag)
	case "config":
		fmt.Print(configTag)
	default:
		handleError("Invalid option")
	}
}
