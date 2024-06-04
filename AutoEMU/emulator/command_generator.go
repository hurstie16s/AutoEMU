package emulator

import (
	"AutoEMU/config"
	"AutoEMU/logger"
	"AutoEMU/machine"
	"fmt"
	"os"
	"os/exec"
	"bufio"
)

func System(machine machine.Machine) (machine.Machine, string, bool) {

	var finalCommand string

	commandRoot := fmt.Sprintf("qemu-system-%s", machine.MachineArch.Arch)

	finalCommand = commandRoot

	// Add cpu tag if cpu is specified
	if machine.CPU != "" {
		finalCommand = fmt.Sprintf("%s -cpu %s", finalCommand, machine.CPU)
	}

	// Add memory tag
	// Ensure memory is not set to 0G
	if machine.Memory == 0 {
		machine.Memory = 1
	}
	finalCommand = fmt.Sprintf("%s -m %dG", finalCommand, machine.Memory)

	// Add drive tag (RFS image)
	finalCommand = fmt.Sprintf("%s -drive file=%s", finalCommand, machine.MachineFiles.HardDriveFile)

	// Add device info
	finalCommand = fmt.Sprintf("%s -device e1000,netdev=%s -netdev user,id=%s,hostfwd=tcp::2222-:22", finalCommand, machine.Name, machine.Name)

	// Add kernel
	finalCommand = fmt.Sprintf("%s -kernel %s", finalCommand, machine.MachineFiles.Kernel)

	// Add initrd
	if machine.MachineFiles.Initrd != "" {
		finalCommand = fmt.Sprintf("%s -initrd %s", finalCommand, machine.MachineFiles.Initrd)
	}

	// Add appends
	finalCommand = fmt.Sprintf("%s -nographic -append \"root=LABEL=%s console=ttyS0\"\n\n", finalCommand, machine.Name)

	/*
	var cmd *exec.Cmd = exec.Command(
		fmt.Sprintf("qemu-system-%s", machine.MachineArch.Arch),
		"-machine",
		machine.Board.Board,
		"-cpu",
		machine.CPU,
		"-m",
		fmt.Sprintf("%dG", machine.Memory),
		"-drive",
		fmt.Sprintf("file=%s", machine.MachineFiles.HardDriveFile),
		"-device",
		"e1000,netdev=net",
		"-netdev",
		"user,id=net,hostfwd=tcp::2222-:22",
		"-kernel",
		machine.MachineFiles.Kernel,
		"-initrd",
		machine.MachineFiles.Initrd,
		"-nographic",
		"-append",
		"\"root=LABEL=rootfs\"",
	)
	*/

	logger.Success("QEMU Command generated", true)

	shellScript := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		"#!/bin/bash",
		fmt.Sprintf("echo \"Running QEMU command for machine %s\"", machine.Name),
		finalCommand,
	)

	// Write shellScript to .sh file in Machine dir
	shellScriptPath := fmt.Sprintf(
		"%s/FirmwareFiles/%s/%s.sh",
		config.Configuration.WorkSpace,
		machine.Name,
		machine.Name,
	)
	shellFile, err := os.Create(shellScriptPath)
	if err != nil {
		logger.Error("Failed to create shell script for QEMU command, command will still run", true)
		return machine, finalCommand, false
	}
	defer shellFile.Close()
	// Write to file
	
	_, err = shellFile.WriteString(shellScript)
	if err != nil {
		logger.Error("Failed to write QEMU command to shell file , command will still run", true)
		return machine, finalCommand, false
	}

	// Ensure file can be executed
	os.Chmod(shellScriptPath, 0755)

	logger.Success("QEMU command saved to shell file", true)
	machine.MachineFiles.ShellScript = fmt.Sprintf("%s.sh", machine.Name)

	return machine, finalCommand, true
}

func runScript(dir string, script string) {

	logger.Info("Starting script run", true)
	cmd := exec.Command(
		"lxterminal",
		"-e",
		fmt.Sprintf("./%s",script),
	)
	cmd.Dir = dir

	out, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	run := make(chan error)
	scanner := bufio.NewScanner(out)

	go func() {
		logger.Info("Command Running", true)

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}

		// check for errors, kill early if need to

		run <- nil
	}()
	
	err = cmd.Start()
	<-run
	err = cmd.Wait()
	//logger.Warn(err.Error(), true)
}
