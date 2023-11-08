package application

import (
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

// ConfigApplicationDefaultYAML is the config of the default application from yaml file.
type ConfigApplicationDefaultYAML struct {
	// cases config
	Cases struct {
		// cases file path
		FilePath string `yaml:"file_path"`
	} `yaml:"cases"`
	// database config
	Database struct {
		// database address
		Address string `yaml:"address"`
		// database user
		User string `yaml:"user"`
		// database password
		Password string `yaml:"password"`
	} `yaml:"database"`
	// server config
	Server struct {
		// server address
		Address string `yaml:"address"`
	} `yaml:"server"`
}


// ConfigApplicationDefault is the config of the default application.
func NewConfigApplicationDefaultFromYAML(filePath string) (cfg *ConfigApplicationDefault, err error) {
	// config
	var cfgYAML ConfigApplicationDefaultYAML

	// read yaml file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("application: open config file error: %w", err)
	}
	defer file.Close()

	// decode yaml file
	err = yaml.NewDecoder(file).Decode(&cfgYAML)
	if err != nil {
		return nil, fmt.Errorf("application: decode config file error: %w", err)
	}

	// serialize
	cfg = &ConfigApplicationDefault{
		CasesFilePath: cfgYAML.Cases.FilePath,
		Database: &mysql.Config{
			Addr: cfgYAML.Database.Address,
			User: cfgYAML.Database.User,
			Passwd: cfgYAML.Database.Password,
		},
		ServerAddress: cfgYAML.Server.Address,
	}

	return
}