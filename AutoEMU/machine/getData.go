package machine

import (
	"AutoEMU/logger"

	"fmt"
	"strconv"
	"strings"
)

func (machine *Machine) GetMachineType(line string) {
	data := getLineData(line)
	dataCheckArch := strings.ToUpper(data)
	switch {
	case strings.Contains(dataCheckArch, "ARM"):
		fallthrough
	case strings.Contains(dataCheckArch, "AARCH64"):
		machine.MachineArch.Arch = "aarch64"
		machine.MachineArch.ArchVariant = data
		machine.MachineArch.Flag = true
	case strings.Contains(dataCheckArch, "MIPS"):
		// MIPS = big endian
		// MIPSEL = small endian
		var arch string = ""
		if !machine.Endianness.Endian {
			arch = "el"
		}

		var bits string = ""
		if machine.ArchBits == 64 {
			bits = "64"
		}

		machine.MachineArch.Arch = fmt.Sprintf("mips%s%s", bits, arch)
		machine.MachineArch.ArchVariant = data

		if machine.Endianness.Flag && machine.ArchBits == 64 {
			machine.MachineArch.Flag = true
		}
	default:
		logger.Error("Failed to read machine type", true)
		return
	}
}

func (machine *Machine) GetMachineEndian(line string) {
	data := getLineData(line)
	if strings.Contains(data, "big endian") {
		logger.Success("Big Endian Detected", true)
		machine.Endianness.Endian = true
	} else {
		logger.Success("Little Endian Detected", true)
		machine.Endianness.Endian = false
	}
	machine.Endianness.Flag = true
}

func (machine *Machine) GetMachineArchBits(line string) {
	flag := machine.ArchBits == 0 
	data := getLineData(line)
	bytes := strings.Split(data, " ")
	bits, _ := strconv.Atoi(bytes[4])
	bits = bits * 32
	machine.ArchBits = int8(bits)
	msg := fmt.Sprintf("ArchBits set to %d", bits)
	logger.Success(msg, flag)
}

func (machine *Machine) GetKernelPreBuilt() {
}
