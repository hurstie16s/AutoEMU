package machine

import (
	"AutoEMU/logger"
	"AutoEMU/shared"

	"fmt"
	"os/exec"
	"slices"
	"strings"
)

// /home/pi5/Documents/uni/ip/WR802N/WR802N.bin

type Machine struct {
	Name              string `default:""`
	MachineFiles      Files
	FirmwareExtension string `default:""`
	ArchBits          int8   `default:"0"`
	MachineArch       Arch
	Endianness        Endian
	Board             Board  
	CPU               string `default:""`
	CPUs              int    `default:"0"`
	Memory            int    `default:"0"`
	FileSystem        string `default:""`
}

type Files struct {
	FirmwareFile      string `default:""`
	HardDriveFile     string `default:""`
	Kernel            string `default:""`
	Initrd            string `default:""`
	ShellScript       string `default:""`
}

type Endian struct {
	Flag   bool `default:"false"`
	Endian bool `default:"true"`
}

type Board struct {
	Flag bool `default:"false"`
	Board string `default:""`
}

type Arch struct {
	Flag bool `default:"false"`
	Arch string `default:""`
	ArchVariant string `default:""`
}

var SupportedArchitectures = []string{
	"arm",
	"mips",
}

// Struct Functions



func (machine *Machine) PrintSupportedBoards() {
	if !machine.MachineArch.Flag {
		logger.Warn("No architecture set, cannot get supported boards", true)
	} else {
		fmt.Println(shared.Separator)
		cmd := exec.Command(
			fmt.Sprintf("qemu-system-%s",machine.MachineArch.Arch),
			"-M",
			"help",
		)
		out, _ := cmd.Output()
		fmt.Print(string(out))
		fmt.Println(shared.Separator)
	}
}

func (machine *Machine) CheckBoardSupported(board string) bool {
	cmd := exec.Command(
		fmt.Sprintf("qemu-system-%s", machine.MachineArch.Arch),
		"-M",
		"help",
	)
	out, _ := cmd.Output()
	boards := strings.Split(string(out), "\n")[1:]
	var finalBoards = []string{}
	for _, b := range boards {
		finalBoards = append(finalBoards, strings.Split(b, " ")[0])
	}

	return slices.Contains(finalBoards, board)
}



// Helper Functions

func getLineData(line string) string {
	i := strings.Index(line, ":")
	if i == -1 {
		return line
	}
	return strings.TrimSpace(line[i+1:])
}

func PrintSupportedArch() {
	fmt.Println("Develop")
}