package utility

import (
	"fmt"
	"github.com/spf13/viper"
)

// Configuration struct definition.
type Config struct {
	Consumer string `yaml:"consumer_tag"`
	Type string `yaml:"type"`
	QueueName string `yaml:"queue_name"`
	ApiKey string `yaml:"api_key"`
	Client struct{
		Uri string `yaml:"uri, omitempty"`
		Queue struct{
			Durable bool `yaml:"durable, omitempty"`
			AutoDelete bool `yaml:"autodelete, omitempty"`
			Exclusive bool `yaml:"exclusive, omitempty"`
			NoWait bool `yaml:"nowait, omitempty"`
			Args map[string]string `yaml:"args, omitempty"`
		} `yaml:"queue, omitempty"`
		Prefetch struct{
			Count int `yaml:"count, omitempty"`
			Size int `yaml:"size, omitempty"`
			Global bool `yaml:"global, omitempty"`
		} `yaml:"prefetch"`
		Consume struct{
			AutoAck bool `yaml:"autoack, omitempty"`
			Exclusive bool `yaml:"exclusive, omitempty"`
			NoLocal bool `yaml:"nolocal, omitempty"`
			NoWait bool `yaml:"nowait, omitempty"`
			Args map[string]string `yaml:"args, omitempty"`
		} `yaml:"consume"`
	} `yaml:"client"`
}

// Check config file and build configuration array.
func ConfigSetup(consumer string, configDir string) Config {

	// Set default config values.
	defaultConfigSet()

	// YML based configuration.
	viper.SetConfigType("yaml")

	// Set file name.
	viper.SetConfigName(consumer)

	// Set default value for consumer.
	viper.SetDefault("consumer_tag", consumer)

	// If config directory flag is empty we use defaults from user home folder or current folder.
	// Otherwise use provided configuration directory from -config flag.
	if len(configDir) == 0 {
		viper.AddConfigPath("$HOME/.rabbinator")
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(configDir)
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

	// Define configuration variable.
	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	// Check supported queue type.
	if config.Type != "mandrill" && config.Type != "mailchimp" {
		panic(fmt.Errorf("you are using unsupported provider type"))
	}

	// Channel is required.
	if config.QueueName == "" {
		panic(fmt.Errorf("queue name which you want to consume is required to be defined in yaml file. Yaml key: queue_name"))
	}

	return config
}

// Set default values.
func defaultConfigSet() {

	// We set default here, maybe some providers does not require API key.
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
