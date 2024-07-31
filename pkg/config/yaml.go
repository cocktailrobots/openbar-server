package config

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"os"
)

type ReversePinConfig struct {
	Pin         int  `yaml:"pin"`
	ForwardHigh bool `yaml:"forward-high"`
}

type DebugHardwareConfig struct {
	NumPumps int    `yaml:"num-pumps"`
	OutFile  string `yaml:"out-file"`
}

type GpioHardwareConfig struct {
	Pins []int `yaml:"pins"`
}

type SequentHardwareConfig struct {
	ExpectedBoardCount int   `yaml:"expected-board-count"`
	RelayMapping       []int `yaml:"relay-mapping"`
}

type HardwareConfig struct {
	Debug   *DebugHardwareConfig   `yaml:"debug"`
	Gpio    *GpioHardwareConfig    `yaml:"gpio"`
	Sequent *SequentHardwareConfig `yaml:"sequent"`
}

type GpioButtonConfig struct {
	Pins          []int `yaml:"pins"`
	DebounceNanos int64 `yaml:"debounce-duration"`
	ActiveLow     bool  `yaml:"active-low"`
	PullUp        bool  `yaml:"pull-up"`
}

type ButtonConfig struct {
	Gpio *GpioButtonConfig `yaml:"gpio"`
}

type DBConfig struct {
	Host *string `yaml:"host"`
	Port *int    `yaml:"port"`
	User *string `yaml:"user"`
	Pass *string `yaml:"pass"`
}

type ListenerConfig struct {
	Port *int    `yaml:"port"`
	Host *string `yaml:"host"`
}

func (c *ListenerConfig) GetHost() string {
	if c.Host == nil {
		return "0.0.0.0"
	}

	return *c.Host
}

func (c *ListenerConfig) GetPort() int {
	if c.Port == nil {
		return 80
	}

	return *c.Port
}

type Config struct {
	Hardware     *HardwareConfig   `yaml:"hardware"`
	ReversePin   *ReversePinConfig `yaml:"reverse-pin"`
	Buttons      *ButtonConfig     `yaml:"buttons"`
	DB           *DBConfig         `yaml:"db"`
	CocktailsApi *ListenerConfig   `yaml:"cocktails-api"`
	OpenBarApi   *ListenerConfig   `yaml:"openbar-api"`
	MigrationDir string            `yaml:"migration-dir"`
}

func Read(filename string, logger *zap.Logger) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", filename, err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config file %s: %w", filename, err)
	}

	configStr, err := yaml.Marshal(&config)
	if err == nil {
		logger.Info("Read Config", zap.String("config_file", filename), zap.String("config", string(configStr)))
	}

	return &config, nil
}
