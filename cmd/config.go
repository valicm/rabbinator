package cmd

import (
	"fmt"
	"github.com/spf13/viper"
)

// Check config file and build configuration array.
func ConfigSetup(channel string, configFile string) {

	viper.SetDefault("channel", channel)

	// Set default config values.
	defaultConfigSet()

	// YML based configuration.
	viper.SetConfigType("yaml")

	// If config flag is empty we use defaults from user home folder.
	// Otherwise use provided configuration file from -config flag.
	if len(configFile) == 0 {
		viper.SetConfigName(channel)
		viper.AddConfigPath("$HOME/.rabbinator")
	} else {
		viper.SetConfigFile(configFile)
	}

	// Check for configuration files.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic(fmt.Errorf("Config file does not exist: %s \n", err))
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	// Run default configuration.
	var queueType = viper.GetString("type")

	// Check queue type.
	if queueType != "mandrill" && queueType != "mailchimp" {
		panic(fmt.Errorf("not supported queue type or values is non existing"))
	}
}

// Set default values.
func defaultConfigSet() {

	viper.SetDefault("type", "")
	viper.SetDefault("api_key", "")

	// Set rabbitmq settings.
	viper.SetDefault("client.uri", "amqp://guest:guest@foreo.loc:5672")

	viper.SetDefault("client.queue.durable", true)
	viper.SetDefault("client.queue.autodelete", false)
	viper.SetDefault("client.queue.exclusive", false)
	viper.SetDefault("client.queue.nowait", false)

	viper.SetDefault("client.prefetch.count", 1)
	viper.SetDefault("client.consume.size", 0)
	viper.SetDefault("client.consume.global", false)

	viper.SetDefault("client.consume.autoack", false)
	viper.SetDefault("client.consume.exclusive", false)
	viper.SetDefault("client.consume.nolocal", false)
	viper.SetDefault("client.consume.noWait", false)

}
