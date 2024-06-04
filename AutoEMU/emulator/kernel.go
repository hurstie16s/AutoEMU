package emulator

import (
	"AutoEMU/autoEmuErrors"
	"AutoEMU/commands"
	"AutoEMU/config"
	"AutoEMU/logger"
	"AutoEMU/machine"
	"AutoEMU/shared"
	
	"os"
	"slices"
	"strings"
	"bufio"
	"fmt"
	"os/exec"
)

func getKernel(machine machine.Machine) machine.Machine {
	for ok := true; ok; {
		// Check if user wants to build custom kernel
		fmt.Println(shared.Separator)
		fmt.Print(
			"Kernel and RFS options:\n" +
			"1)Build kernel with buildroot\n" +
			"2)Use pre-built kernel for select architecture\n" +
			"3)Provide custom pre-built kernel (Not implemented)\n",
		)
		input := shared.GetInputInteger("Invalid Option", "Invalid Option")
		switch input {
		case 1:
			machine, _ = buildKernel(machine)
			ok = false
		case 2:
			machine, ok = getKernelPreBuilt(machine)
		case 3:
			logger.Error("Functionality Not Implemented", true)
		default:
			logger.Error("Invalid Option", true)
		}
	}

	return machine
}

func buildKernel(machine machine.Machine) (machine.Machine, error) {

	//var out []byte

	// Print important information about machine
	fmt.Println(shared.Separator)
	fmt.Printf("Machine Information for %s\n", machine.Name)
	fmt.Printf("Architecture: %s\n", machine.MachineArch.Arch)
	fmt.Printf("Architecture Variant: %s\n", machine.MachineArch.ArchVariant)
	if machine.Endianness.Endian {
		fmt.Println("Endian: Big Endian")
	} else {
		fmt.Println("Endain: Little Endian")
	}
	fmt.Printf("%d bit architecture\n", machine.ArchBits)
	if machine.CPU != "" {
		fmt.Printf("CPU: %s\n", machine.CPU)
	}
	fmt.Printf("Board: %s\nThis is important for the kernel defconfig name\n", machine.Board.Board)
	fmt.Printf("Root File System type: %s\n", machine.FileSystem)
	fmt.Println(shared.Separator)
	fmt.Print("Press Enter to start buildroot")
	fmt.Scanln()

	dir := fmt.Sprintf(
		"%s/BuildRoot",
		config.Configuration.WorkSpace,
	)
	
	// run make distclean
	cmd := exec.Command(
		"make",
		"distclean",
	)
	cmd.Dir = dir
	cmd.Run()
	logger.Info("Previous configs cleared", true)
	logger.Info("Running nconfig", true)

	// run make nconfig
	cmd = exec.Command(
		"lxterminal",
		"-e",
		"make",
		"nconfig",
	)
	cmd.Dir = dir
	test, _ := cmd.Output()
	fmt.Println(string(test))

	fmt.Print("Press Enter to start buildroot : WARNING: DO NOT DO THIS UNTIL YOU HAVE SAVED YOUR KERNEL CONFIG")
	fmt.Scanln()
	logger.Info("Building kernel and RFS", true)
	logger.Warn("This process may take a long time, do not terminate the application", true)

	cmd = exec.Command("make")
	cmd.Dir = dir
	
	out, _ := cmd.StdoutPipe()

	cmd.Stderr = cmd.Stdout

	run := make(chan error)

	scanner := bufio.NewScanner(out)

	go func() {

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

	// Move kernel file and RFS into machine dir
	buildRootOutDir := fmt.Sprintf(
		"%s/BuildRoot/output/images",
		config.Configuration.WorkSpace,
	)
	files, err := os.ReadDir(buildRootOutDir)
	if err != nil {
		return machine, &autoEmuErrors.KernelOutputReadError{Err: err}
	}
	if len(files) < 2 {
		return machine, &autoEmuErrors.KernelBuildError{
			Err: fmt.Errorf("image files for kernel and RFS not present in buildroot output"),
		}
	}
	var imagesFound bool = false
	// Get kernel file
	for _, file := range files {
		if !strings.Contains(file.Name(), ".") {
			logger.Info("Found kernel File", true)
			// Copy file
			newPath := fmt.Sprintf(
				"FirmwareFiles/%s/%s",
				machine.Name,
				machine.Name,
			)
			currentFile := fmt.Sprintf(
				"BuildRoot/output/images/%s",
				file.Name(),
			)
			commands.CopyFile(currentFile, newPath)
			machine.MachineFiles.Kernel = machine.Name

			imagesFound = true

			break
		}
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	if slices.Contains(fileNames, fmt.Sprintf("rootfs.%s", "ext2")) {
		newPath := fmt.Sprintf(
			"FirmwareFiles/%s/%s.%s",
			machine.Name,
			machine.Name,
			"ext2",
		)
		currentPath := fmt.Sprintf(
			"BuildRoot/output/images/rootfs.%s",
			"ext2",
		)
		commands.CopyFile(currentPath, newPath)
		machine.MachineFiles.HardDriveFile = fmt.Sprintf("%s.%s", machine.Name, "ext2")
		imagesFound = imagesFound && true
	} else if slices.Contains(fileNames, "rootfs.tar") {
		logger.Warn("No RFS avaliale for machine filesystem", true)
		logger.Info("Building RFS from tar file", true)
		// Build RFS from tar
		logger.Error("Not Implemented", true)
		imagesFound = imagesFound && true
	} else {
		imagesFound = false
	}

	if !imagesFound {
		return machine, &autoEmuErrors.KernelBuildError{
			Err: fmt.Errorf("kernel and RFS image files not present in workspace"),
		}
	}

	return machine, nil
}

func getKernelPreBuilt(machine machine.Machine) (machine.Machine, bool) {

	switch machine.MachineArch.Arch {
	case "mipsel":
		logger.Info("Getting Pre-built kernel, RFS and Initrd for mispel architecture", true)
		// Copy Kernel into dir
		commands.CopyFile(
			"Kernels/mipselImages/mipselKernel",
			fmt.Sprintf("FirmwareFiles/%s", machine.Name),
		)
		// Copy RFS into dir
		commands.CopyFile(
			"Kernels/mipselImages/mipselRFS.squashfs",
			fmt.Sprintf("FirmwareFiles/%s", machine.Name),
		)
		// Copy Initrd into dir
		commands.CopyFile(
			"Initrds/mipselInitrd",
			fmt.Sprintf("FirmwareFiles/%s", machine.Name),
		)
	default:
		logger.Error("No prebuilt kernel avaliable for your architecture", true)
		return machine, true
	}

	machine.MachineFiles.HardDriveFile = fmt.Sprintf("%sRFS.squashfs", machine.MachineArch.Arch)
	machine.MachineFiles.Kernel = fmt.Sprintf("%sKernel", machine.MachineArch.Arch)
	machine.MachineFiles.Initrd = fmt.Sprintf("%sInitrd", machine.MachineArch.Arch)

	return machine, false
}