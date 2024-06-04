package config

//Imports
import (
	"AutoEMU/logger"
	"errors"
	"fmt"
	"os"
	"strings"
	"AutoEMU/shared"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Configured bool `default:"true"`
	WorkSpace  string `default:""`
	LogDir     string `default:""`
}

var fields = []string{
	"workspace directory",
	"logs directory",
}
var configVal Config
var Configuration *Config = &configVal
var err error
var flag bool
var input string = ""

func loadViper() error {

	viper.AddConfigPath(pwd())
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	return viper.ReadInConfig()
}

func LoadConfig(ignoreConfig bool) error {

	Configuration = &Config{
		Configured: false,
		WorkSpace: "",
		LogDir: "",
	}

	if ignoreConfig {
		return nil
	}

	err = loadViper()
	if err != nil {

		createConfigFile()

		// Reattempt to loadViper
		viper.Reset()
		err = loadViper()
		if err != nil {
			return fmt.Errorf("cannot load or create config: %v", err)
		}
	}


	err = viper.Unmarshal(&Configuration)
	if err != nil {
		return fmt.Errorf("could not unmarshal config file: %v", err)
	}

	if flag {
		// Build config file with user options
		buildConfig()
		writeConfig()
	}
	
	return nil
}

func createConfigFile() {
	flag = true
	var configFile *os.File
	fmt.Printf("%sNo Config Created, Creating Config\n", logger.ConfigMessage)
	configFile, err = os.Create("config.yaml")
	if err != nil {
		panic(err)
	}
	configFile.Close()
}

func writeConfig() {
	configFile, _ := os.OpenFile("config.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer configFile.Close()
	encoder := yaml.NewEncoder(configFile)
	_ = encoder.Encode(Configuration)
	encoder.Close()
}

func buildConfig() {
	var len int

	fmt.Println(shared.Separator)
	fmt.Println("Building configuration for AutoEMU")
	fmt.Println(shared.Separator)

	// Get workspace
	fmt.Println("Enter path to WorkSpace directory or leave blank if WorkSpace is in this directory")
	fmt.Print(">>>")
	len, _ = fmt.Scanln(&input)
	var workSpace string
	if len == 0 {
		workSpace = fmt.Sprintf("%s/WorkSpace", pwd())
	} else {
		workSpace = input
		// Verify that path exists
	}
	Configuration.SetConfigWorkSpace(workSpace)
	
	// Get Logs DIR
	fmt.Println("Enter path to a directory to create a Log directory or leave blank to select this directory")
	fmt.Print(">>>")
	len, _ = fmt.Scanln(&input)
	var logsDir string
	if len == 0 {
		input = pwd()
	}
	logsDir = fmt.Sprintf("%s/Logs", input)
	// Create the dir if not exist
	if _, err := os.Stat(logsDir); errors.Is(err, os.ErrNotExist) {
		_ = os.Mkdir(logsDir, os.ModePerm)
	}
	Configuration.SetConfigLogDIR(logsDir)

	Configuration.SetConfigFlag(true)

	fmt.Println("===================================================================================================")
	fmt.Println("Configuration built")
	fmt.Println("===================================================================================================")

}

func BuildConfig(workSpace string, logs string) {

	if workSpace == "-pwd" {
		workSpace = fmt.Sprintf("%s/WorkSpace", pwd())
	}

	if logs == "-pwd" {
		logs = pwd()
	}
	logs = fmt.Sprintf("%s/Logs", logs)

	// Create the dir if not exist
	if _, err := os.Stat(logs); errors.Is(err, os.ErrNotExist) {
		_ = os.Mkdir(logs, os.ModePerm)
	}

	Configuration.SetConfigWorkSpace(workSpace)
	Configuration.SetConfigLogDIR(logs)
	Configuration.SetConfigFlag(true)

	writeConfig()
}

func (config *Config ) SetConfigWorkSpace(workSpace string) {
	config.WorkSpace = workSpace
}

func (config *Config) SetConfigFlag(flag bool) {
	config.Configured = flag
}

func (config *Config) SetConfigLogDIR(dir string) {
	config.LogDir = dir
}

func pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		// TODO Handle error
		return ""
	}
	return dir
}

func (config *Config) UpdateConfig(field int) {
	msg := fmt.Sprintf("Updating %s", fields[field-1])
	logger.Config(msg, true)
	var update bool = false
	// Get value
	fmt.Printf("Enter value for %s or -m to return to the menu\n", fields[field-1])
	for {
		input := shared.GetInput("Please provide an input or -m to return to menu")
		_, err = os.Stat(input)
		if strings.ToLower(input) == "-m" {
			logger.Config("Canceling configuration update", true)
			logger.Info("Returning to menu", false)
			update = false
			break
		} else if os.IsNotExist(err) {
			logger.Error("Could not find resource", true)
		} else {
			update = true
			break
		}
	}
	if update {
		switch field {
		case 1:
			Configuration.SetConfigWorkSpace(input)
			logger.Config("Workspace directory updated", true)
		case 2:
			Configuration.SetConfigLogDIR(input)
			logger.Config("Logs directory updated", true)
		}
		writeConfig()
	}
	// Return to menu
}

func CheckWorkSpace() {
	// Check FirmwareFiles
	if !exists(fmt.Sprintf("%s/FirmwareFiles", Configuration.WorkSpace)) {
		os.Mkdir(fmt.Sprintf("%s/FirmwareFiles", Configuration.WorkSpace), os.ModePerm)
	}
	// Check ZippedMachines
	if !exists(fmt.Sprintf("%s/ZippedMachines", Configuration.WorkSpace)) {
		os.Mkdir(fmt.Sprintf("%s/ZippedMachines", Configuration.WorkSpace), os.ModePerm)
	}
}

func exists(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}