package emulator

import (
	"AutoEMU/autoEmuErrors"
	"AutoEMU/logger"
	"AutoEMU/shared"
	"AutoEMU/machine"

	"errors"
	"fmt"
	"slices"
	"strings"
)

var err error = errors.New("TMP")
var input string

func checkAnalysis(machineConfig machine.Machine) machine.Machine {

	// Check for endianness
	if !machineConfig.Endianness.Flag {
		logger.Warn("Machine endianness could not be automatically detected", true)

		for machineConfig.Endianness.Flag {

			input = strings.ToUpper(
				shared.GetInputWithPrompt(
					"Enter machine endianness (big/ little)",
					"No endian entered",
				),
			)
			switch input {
				case "L":
					machineConfig.Endianness.Endian = false
					machineConfig.Endianness.Flag = true
				case "B":
					machineConfig.Endianness.Endian = true
					machineConfig.Endianness.Flag = true
				default:
					logger.Error("Invalid option, enter big or little", true)
			}

		}
	}

	// Check for bits
	if machineConfig.ArchBits == 0 {
		logger.Warn("AutoEMU could not determine a 32 or 64 bit architecture", true)

		var bits int8
		for {
			bits = int8(shared.GetInputWithPromptInteger(
				"Is the machine 32 or 64 bit architecture, enter 32 or 64",
				"No architecture entered",
				"Invalid architecture",
			))
			if bits == 32 || bits == 64 {
				break
			}
			logger.Error("Invalid architecture", true)
		}
	}

	// Check if architecture has been detected
	if !machineConfig.MachineArch.Flag {
		logger.Warn("Machine architecture could not be automatically detected", true)

		for ok := true; ok; ok = err != nil {
			
			input = shared.GetInputWithPrompt(
				"Enter machine architectutre or -list to list supported architectures",
				"No architecture entered",
			)
			
			if input == "-list" {
				machine.PrintSupportedArch()
				machineConfig.PrintSupportedBoards()
			} else {
				msg := fmt.Sprintf("Architecture Selected: %v", input)
				logger.Info(msg, false)

				err = checkArchSupported(input)
				if err != nil {
					logger.Error(err.Error(), true)
				} else {
					msg := fmt.Sprintf("Architecture %s accepted", input)
					logger.Success(msg, true)
					machineConfig.GetMachineType(input)
				}
			}
		}
	}

	// Check if board is present
	if !machineConfig.Board.Flag {
		
		for {

			var input string
			fmt.Println("Enter machine board, -list to list supported boards or leave blank if you do not wish to specify a board")
			fmt.Print(">>>")
			len, _ := fmt.Scanln(&input)
			
			if len == 0 {
				break
			} else if input == "-list" {
				machineConfig.PrintSupportedBoards()
			} else {
				if machineConfig.CheckBoardSupported(input) {
					machineConfig.Board.Board = input
					machineConfig.Board.Flag = true
					break
				} else {
					logger.Error(fmt.Sprintf("%s is not a supported board", input), true)
				}
			}

		}

	}

	return machineConfig
}

func checkArchSupported(input string) error {
	var archSupported bool = slices.Contains(machine.SupportedArchitectures, input)
	if !archSupported {
		return &autoEmuErrors.ArchitectureUnsupportedError{Err: errors.New(input)}
	}
	return nil
}