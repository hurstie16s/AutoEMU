package terminal

import (
	"AutoEMU/config"
	"AutoEMU/logger"
	
	"fmt"
	"os"
	"strings"
)

func ProcessArgs(args []string) {
	for ok := true; ok; ok = len(args) != 0 {
		switch args[0] {
		case "-h":
			if len(args) == 1 || strings.HasPrefix(args[1], "-") {
				printFullHelp()
				args = args[1:]
			} else {
				printByTag(args[1])
				args = args[2:]
			}
		default:
			handleError(fmt.Sprintf("Unknown tag: %s", args[0]))
		}
	}
	os.Exit(0)
}

func ProcessConfig(args []string) []string {
	if len(args) < 3 {
		errorMessage := fmt.Sprintf("Expected 2 arguments, found %d", len(args)-1)
		handleError(errorMessage)
	}
	config.BuildConfig(args[1], args[2])
	// Remove elements from args
	args = args[3:]
	fmt.Printf("%sConfiguration built\n", logger.SuccessMessage)
	return args
}

func handleError(msg string) {
	fmt.Printf("%s%s\n", logger.FatalMessage, msg)
	fmt.Printf("%sRun ./AutoEMU -h to view the help page\n", logger.InfoMessage)
	os.Exit(1)
}