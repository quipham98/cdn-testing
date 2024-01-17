package config

import (
	"github.com/trustwallet/go-libs/config/viper"
	"path/filepath"
)

type (
	Config struct {
		App                App                `mapstructure:"app"`
		ValidatorsSettings ValidatorsSettings `mapstructure:"validators_settings"`
	}

	App struct {
		LogLevel string `mapstructure:"log_level"`
	}

	ValidatorsSettings struct {
		RootFolder      RootFolder      `mapstructure:"root_folder"`
		ChainFolder     ChainFolder     `mapstructure:"chain_folder"`
		AssetFolder     AssetFolder     `mapstructure:"asset_folder"`
		ChainInfoFolder ChainInfoFolder `mapstructure:"chain_info_folder"`
		Image           Image           `mapstructure:"image"`
	}

	Image struct {
		MaxH int `mapstructure:"max_h"`
		MaxW int `mapstructure:"max_w"`
		MinH int `mapstructure:"min_h"`
		MinW int `mapstructure:"min_w"`
	}
)

// Default is a configuration instance.
var Default = Config{} //nolint:gochecknoglobals // config must be global

// SetConfig reads a config file and returs an initialized config instance.
func SetConfig(confPath string) error {
	confPath, err := filepath.Abs(confPath)
	if err != nil {
		return err
	}

	viper.Load(confPath, &Default)

	return nil
}
