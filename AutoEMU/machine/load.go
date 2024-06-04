package machine

import (
	//"AutoEMU/config"
	"fmt"

	"github.com/spf13/viper"
)

var machineVal Machine
var machine *Machine = &machineVal

func Load(machineConfig string) (Machine, error) {

	// Load viper
	viper.AddConfigPath(
		fmt.Sprintf(
			"WorkSpace/FirmwareFiles/%s",
			machineConfig,
		),
	)
	viper.SetConfigName(machineConfig)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return machineVal, err
	}

	err = viper.Unmarshal(machine)
	if err != nil {
		return machineVal, err
	}

	return machineVal, nil
}