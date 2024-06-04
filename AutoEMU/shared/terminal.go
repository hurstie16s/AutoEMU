package shared

import (
	"AutoEMU/logger"
	"fmt"
	"strconv"
	"strings"
)

const Separator string = "==================================================================================================="

func GetInput(errorMessage string) string {
	var input string
	var err error
	for ok := true; ok; ok = err != nil {
		fmt.Print(">>>")
		_, err = fmt.Scanln(&input)
		if err != nil {
			logger.Error(errorMessage, true)
		}
	}
	return input
}

func GetInputWithPrompt(prompt string, errorMessage string) string {
	var input string
	var err error
	for ok := true; ok; ok = err != nil {
		fmt.Println(prompt)
		fmt.Print(">>>")
		_, err = fmt.Scanln(&input)
		if err != nil {
			logger.Error(errorMessage, true)
		}
	}
	return input
}

func GetInputWithPromptInteger(prompt string, errorMessageEmptyEntry string, errorMessageNoInt string) int64 {
	var bits int64
	var input string
	var err error
	for ok := true; ok; ok = err != nil {
		input = GetInputWithPrompt(prompt, errorMessageEmptyEntry)
		// Attempt to convert to int8
		bits, err = strconv.ParseInt(input, 10, 8)
		if err != nil {
			logger.Error(errorMessageNoInt, true)
		}
	}
	return bits
}

func GetInputInteger(errorMessageEmptyEntry string, errorMessageNoInt string) int64 {
	var bits int64
	var input string
	var err error
	for ok := true; ok; ok = err != nil {
		input = GetInput(errorMessageEmptyEntry)
		// Attempt to convert to int8
		bits, err = strconv.ParseInt(input, 10, 8)
		if err != nil {
			logger.Error(errorMessageNoInt, true)
		}
	}
	return bits
}

func GetInputBinary(prompt string, opt1 string, opt2 string) bool {
	var input string
	for {
		input = GetInputWithPrompt(prompt, "Invalid Option")
		if strings.EqualFold(input, opt1) {
			return true
		} else if strings.EqualFold(input, opt2) {
			return false
		} else {
			logger.Error("Invalid Option", true)
		}
	}
}