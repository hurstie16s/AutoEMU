package commands

import (
	"AutoEMU/config"
	"testing"
)

func TestCreateDIRFromWorkSpace(t *testing.T) {
	config.Configuration.SetConfigWorkSpace("/home/pi5/Documents/uni/ip/comp3200-part-3-ip/AutoEMU/testsWorkSpace")
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "3: testToFail",
			args: args{
				path: "FirmwareFiles/takenMachineName",
			},
			wantErr: true,
		},
		{
			name: "2: testToPass",
			args: args{
				path: "FirmwareFiles/newMachineName",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateDIRFromWorkSpace(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("CreateDIRFromWorkSpace() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("Test %s passed", tt.name)
			}
		})
	}
}
