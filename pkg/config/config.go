package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// InitConfig ...
func InitConfig(prefix string) {
	viper.SetConfigType("yaml")

	configDir := os.Getenv("CONFIG_DIR")
	if configDir != "" {
		viper.AddConfigPath(configDir)
	}
	name := os.Getenv("GO_ENV")
	if name == "" {
		name = "development"
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")
	viper.AddConfigPath("/etc/" + prefix)

	viper.AutomaticEnv()
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName(name)
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(err)
	}

	if viper.GetBool("watchConfigChange") {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Printf("Config file changed: %s, reloading ...\n", e.Name)
			if err := viper.ReadInConfig(); err == nil {
				fmt.Printf("Reload config file success: %s\n", viper.ConfigFileUsed())
			} else {
				fmt.Printf("Reload config file error %v\n", err)
			}
		})
	}
}

type ConfigUpdateWatcher interface {
	OnConfigUpdate(func(cfg *viper.Viper))
}
