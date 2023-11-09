package application

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ConfigApplicationDefaultYAML is the config of the default application from yaml file.
type ConfigApplicationDefaultYAML struct {
	// server config
	Server struct {
		// server address
		Address string `yaml:"address"`
	} `yaml:"server"`
	// database config
	Database struct {
		// database address
		Address string `yaml:"address"`
		// database user
		User string `yaml:"user"`
		// database password
		Password string `yaml:"password"`
		// database name
		Name string `yaml:"name"`
	} `yaml:"database"`
	// cases config
	Cases struct {
		Reader struct {
			FilePath string `yaml:"file_path"`
			BatchSize int `yaml:"batch_size"`
		} `yaml:"reader"`
		Reporter struct {
			ExcludedHeaders []string `yaml:"excluded_headers"`
		} `yaml:"reporter"`
	} `yaml:"cases"`
}


// ConfigApplicationDefault is the config of the default application.
func NewConfigApplicationDefaultFromYAML(filePath string) (cfg *Config, err error) {
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
	cfg = &Config{
		Server: ServerConfig{
			Address: cfgYAML.Server.Address,
		},
		Database: &DatabaseConfig{
			Address: cfgYAML.Database.Address,
			User: cfgYAML.Database.User,
			Password: cfgYAML.Database.Password,
			Name: cfgYAML.Database.Name,
		},
		Cases: CasesConfig{
			Reader: struct {
				FilePath string
				BatchSize int
			}{
				FilePath: cfgYAML.Cases.Reader.FilePath,
				BatchSize: cfgYAML.Cases.Reader.BatchSize,
			},
			Reporter: struct {
				ExcludedHeaders []string
			}{
				ExcludedHeaders: cfgYAML.Cases.Reporter.ExcludedHeaders,
			},
		},
	}
	return
}