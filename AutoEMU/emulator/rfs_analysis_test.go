package emulator

import (
	"AutoEMU/machine"
	"AutoEMU/config"
	"testing"
)

func TestGetArchInfo(t *testing.T) {
	config.Configuration.SetConfigWorkSpace("/home/pi5/Documents/uni/ip/comp3200-part-3-ip/AutoEMU/WorkSpace")
	type args struct {
		machineConfig machine.Machine
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 1: WR802N - RFS Analysis",
			args: args{
				machine.Machine{
					Name: "WR802N",
					FirmwareExtension: "bin",
					MachineFiles: machine.Files{
						FirmwareFile: "WR802N.bin",
					},
				},
			},
		},
		{
			name: "Test 2: CPE510 - RFS Analysis",
			args: args{
				machine.Machine{
					Name: "CPE510",
					FirmwareExtension: "bin",
					MachineFiles: machine.Files{
						FirmwareFile: "CPE510.bin",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetArchInfo(tt.args.machineConfig)
			if !got.Endianness.Flag {
				t.Error("Failed to determine machine endianess")
				return
			}
			if got.ArchBits == 0 {
				t.Error("Failed to determine 32 or 64 bit architecture")
				return
			}
			if got.MachineArch.Arch == "" || got.MachineArch.ArchVariant == "" {
				t.Error("Failed to determine machine architecture")
			}
		})
	}
}
