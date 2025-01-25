package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func LoadConfig(configFilePath string) (*Config, error) {
	var config = &Config{}

	if _, err := os.Stat(configFilePath); err != nil {
		log.Warn(fmt.Sprintf("Config file '%s' not found, using default config", configFilePath))
		return DefaultConfig(), nil
	}

	configFileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(configFilePath)
	switch strings.ToLower(ext) {
	case ".json":
		err = json.Unmarshal(configFileContent, config)
		if err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		err = yaml.Unmarshal(configFileContent, config)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("couldn't process the file %s", configFilePath)
	}

	// include the default config if the flag is set
	if config.IncludeDefault {
		defaultConfig := DefaultConfig()
		config.Rules = append(config.Rules, defaultConfig.Rules...)
		config.IgnorePaths = append(config.IgnorePaths, defaultConfig.IgnorePaths...)
		config.IgnoreExtensions = append(config.IgnoreExtensions, defaultConfig.IgnoreExtensions...)
		config.Exceptions = append(config.Exceptions, defaultConfig.Exceptions...)
	}

	return config, nil
}
