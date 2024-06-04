package emulator

import (
	"AutoEMU/commands"
	"AutoEMU/config"
	"AutoEMU/logger"
	"AutoEMU/machine"

	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var workingDir string
var machineTMP machine.Machine

func GetArchInfo(machineConfig machine.Machine) machine.Machine {
	machineTMP = machineConfig

	// Set Working Dir
	workingDir = fmt.Sprintf(
		"%s/FirmwareFiles/%s",
		config.Configuration.WorkSpace,
		machineTMP.Name,
	)

	// Make elfs dir
	cmd := exec.Command("mkdir", "elfs")
	cmd.Dir = workingDir
	cmd.Run()

	logger.Info("ELFS dir built", false)

	// Get All ELF executables into elfs dir
	dirRead, err := getFilesRFS()
	if err != nil {
		return machineTMP
	}

	// Remove any directories from dirRead
	var elfs []fs.DirEntry
	for _, entry := range dirRead {
		if !entry.IsDir() {
			elfs = append(elfs, entry)
		}
	}
	// Run readelf -A on each elf file until
	logger.Info("Starting analysis of ELF executables", true)
	for _, elf := range elfs {
		analyseELF(elf.Name())
	}
	if machineTMP.MachineArch.Flag {
		logger.Success(fmt.Sprintf("Machine Type set to: %s", machineTMP.MachineArch.Arch), true)
	}

	return machineTMP
}

func getFilesRFS() ([]fs.DirEntry, error) {

	logger.Info("Start searching for ELF executables", true)

	var dir string = fmt.Sprintf(
		"%s/FirmwareFiles/%s/fmk/rootfs",
		config.Configuration.WorkSpace,
		machineTMP.Name,
	)

	filepath.WalkDir(dir, walk)

	// Get list of files collected
	files, err := os.ReadDir(fmt.Sprintf("%s/elfs", workingDir))
	var length = len(files)
	var msg = fmt.Sprintf("%d ELF executables found", length)
	logger.Info(msg, true)

	if err != nil {
		logger.Warn("Failed to collect ELF files for architecture analysis", true)
		return nil, err
	}
	return files, nil
}

func walk(file string, dir fs.DirEntry, err error) error {

	if err != nil {
		return err
	}
	if !dir.IsDir() {
		checkELF(file)
	}
	return nil
}

func checkELF(filePath string) error {

	// Remove WorkSpace prefix from file path
	var fileStripped string = filePath[len(workingDir)+1:]

	cmd := exec.Command(
		"file",
		fileStripped,
	)
	cmd.Dir = workingDir
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("handle Error")
		os.Exit(1)
	}

	// Remove file name from output incase ELF tag is present in path
	var analysis string = string(out[len(fileStripped):])
	// Check if output contains ELF tag
	if strings.Contains(analysis, "ELF") {
		// Move file into elfs dir
		workSpaceFilePath := fmt.Sprintf("FirmwareFiles/%s", machineTMP.Name)
		err := commands.CopyFile(
			fmt.Sprintf("%s/%s", workSpaceFilePath, fileStripped),
			fmt.Sprintf("%s/elfs", workSpaceFilePath),
		)

		// Add name of elf to ELFS array in machineConfig

		if err == nil {
			return nil
		}
		return nil // Change to return an error
	}
	return nil
}

func analyseELF(elfName string) error {

	msg := fmt.Sprintf("Analysing ELF File %s", elfName)
	logger.Info(msg, false)

	/*
		cmd := exec.Command(
			"readelf",
			"-A",
			fmt.Sprintf("elfs/%s", elfName),
		)

		cmd := exec.Command(
			"objdump",
			"-f",
			elfName,
		)
	*/

	// Get file headers
	cmd := exec.Command(
		"readelf",
		"-h",
		elfName,
	)
	cmd.Dir = fmt.Sprintf("%s/elfs", workingDir)

	out, _ := cmd.StdoutPipe()

	cmd.Stderr = cmd.Stdout

	msg = fmt.Sprintf("Read analysis of %s", elfName)
	logger.Info(msg, false)
	complete := make(chan struct{})

	scanner := bufio.NewScanner(out)

	go func() {

		// Read command output line by line
		for scanner.Scan() {
			line := scanner.Text()
			logger.Write(line, false)

			switch {
			case strings.Contains(line, "Magic:"):
				if machineTMP.ArchBits == 0 {
					logger.Info("Getting Architecture Bits", true)
					machineTMP.GetMachineArchBits(line)
				} else if machineTMP.ArchBits == 32 {
					machineTMP.GetMachineArchBits(line)
				}
			case strings.Contains(line, "Data:"):
				if !machineTMP.Endianness.Flag {
					logger.Info("Getting Endian", true)
					machineTMP.GetMachineEndian(line)
				}
			case strings.Contains(line, "Machine:"):
				if !machineTMP.MachineArch.Flag {
					logger.Info("Getting machine type", true)
					machineTMP.GetMachineType(line)
				}
			}
		}

		if machineTMP.MachineArch.Arch != "" {
			machineTMP.MachineArch.Flag = true
		}

		complete <- struct{}{}
	}()

	_ = cmd.Start()

	<-complete

	_ = cmd.Wait()

	return nil
}

func analyseConfigLog(machineConfig machine.Machine) machine.Machine {

	logger.Info("Analysing config log", true)

	// Log file
	logFile := fmt.Sprintf(
		"%s/FirmwareFiles/%s/fmk/logs/config.log",
		config.Configuration.WorkSpace,
		machineConfig.Name,
	)

	file, err := os.Open(logFile)
	if err != nil {
		return machineConfig
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "FS_TYPE") {
			machineConfig.FileSystem = line[9 : len(line)-1]
			logger.Success(fmt.Sprintf("File system type found: %s", machineConfig.FileSystem), true)
		} else if strings.HasPrefix(line, "ENDIANESS") {
			machineConfig.Endianness.Endian = line[11 : len(line)-1] != "-le"
			machineConfig.Endianness.Flag = true
			logger.Success("Machine Endianess set", true)
		}
	}

	return machineConfig
}
