package emulator

import (
	"AutoEMU/config"
	"AutoEMU/machine"
	"reflect"
	"testing"
)

func TestExtractFirmware(t *testing.T) {

	config.Configuration.WorkSpace = "/home/pi5/Documents/uni/ip/comp3200-part-3-ip/AutoEMU/WorkSpace"

	type args struct {
		pathToFile string
		name       string
	}
	tests := []struct {
		name    string
		args    args
		want    machine.Machine
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test 1: WR802N Firmware Extraction",
			args: args{
				pathToFile: "/home/pi5/Documents/uni/ip/FirmwareZips/WR802N.bin",
				name: "WR802N",
			},
			want: machine.Machine{
				Name: "WR802N",
				FirmwareExtension: "bin",
				MachineFiles: machine.Files{
					FirmwareFile: "WR802N.bin",
				},
			},
			wantErr: false,
		},
		{
			name: "Test 2: TL-WR702N Firmware Extraction",
			args: args{
				pathToFile: "/home/pi5/Documents/uni/ip/FirmwareZips/TL-WR702N.bin",
				name: "TL-WR702N",
			},
			want: machine.Machine{
				Name: "TL-WR702N",
				FirmwareExtension: "bin",
				MachineFiles: machine.Files{
					FirmwareFile: "TL-WR702N.bin",
				},
			},
			wantErr: true,
		},
		{
			name: "Test 3: GT-AX11000 Firmware Extraction",
			args: args{
				pathToFile: "/home/pi5/Documents/uni/ip/FirmwareZips/GT-AX11000.w",
				name: "GT-AX11000",
			},
			want: machine.Machine{
				Name: "GT-AX11000",
				FirmwareExtension: "w",
				MachineFiles: machine.Files{
					FirmwareFile: "GT-AX11000.w",
				},
			},
			wantErr: true,
		},
		{
			name: "Test 4: CPE510-V3-20 Frimware Extraction",
			args: args{
				pathToFile: "/home/pi5/Documents/uni/ip/FirmwareZips/CPE510-V3-20.bin",
				name: "CPE510",
			},
			want: machine.Machine{
				Name: "CPE510",
				FirmwareExtension: "bin",
				MachineFiles: machine.Files{
					FirmwareFile: "CPE510.bin",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractFirmware(tt.args.pathToFile, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFirmware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractFirmware() = %v, want %v", got, tt.want)
			}
		})
	
	}
}
