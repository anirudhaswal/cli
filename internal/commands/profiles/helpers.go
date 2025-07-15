package profiles

import (
	"bufio"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Profile struct {
	BaseUrl      string `yaml:"SUPRSEND_BASE_URL"`
	MgmntUrl     string `yaml:"SUPRSEND_MGMNT_URL"`
	ServiceToken string `yaml:"SUPRSEND_SERVICE_TOKEN"`
}

type Config struct {
	ActiveProfile string             `yaml:"active_profile"`
	Profiles      map[string]Profile `yaml:"profiles"`
}

func EnsureConfig(path string) (*Config, string, error) {
	var configPath string
	if path != "" {
		configPath = path
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.WithError(err).Error("Could not get user home directory")
			return nil, "", err
		}
		configPath = filepath.Join(homeDir, ".suprsend.yaml")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Warnf("No config found at %s", configPath)
		log.Info("Would you like to create a default config? (Y/n): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		response := scanner.Text()
		if response != "y" && response != "Y" {
			log.Error("Config file is required to proceed.")
			return nil, configPath, err
		}

		defaultCfg := &Config{
			ActiveProfile: "default",
			Profiles: map[string]Profile{
				"default": {
					BaseUrl:      "https://hub.suprsend.com",
					MgmntUrl:     "https://api.suprsend.com",
					ServiceToken: "",
				},
			},
		}
		if err := SaveConfig(defaultCfg, configPath); err != nil {
			log.WithError(err).Error("Failed to create default config")
			return nil, configPath, err
		}
		log.Infof("Created default config at %s", configPath)
		return defaultCfg, configPath, nil
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.WithError(err).Error("Failed to load config file")
		return nil, configPath, err
	}
	return cfg, configPath, nil
}

func SaveConfig(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.WithError(err).Errorf("Failed to read config file: %s", path)
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.WithError(err).Error("Failed to parse YAML config")
		return nil, err
	}
	return &cfg, nil
}
