package utility

import (
	"github.com/spf13/viper"
)

// Configuration struct definition.
type Config struct {
	Type      string `yaml:"type"`
	QueueName string `yaml:"queueName"`
	ApiKey    string `yaml:"apiKey"`
	Consumer  string `yaml:"consumerTag"`
	Client    struct {
		Uri   string `yaml:"uri"`
		Queue struct {
			Durable    bool `yaml:"durable, omitempty"`
			AutoDelete bool `yaml:"autodelete, omitempty"`
			Exclusive  bool `yaml:"exclusive, omitempty"`
			NoWait     bool `yaml:"nowait, omitempty"`
		} `yaml:"queue"`
		Prefetch struct {
			Count  int  `yaml:"count, omitempty"`
			Size   int  `yaml:"size, omitempty"`
			Global bool `yaml:"global, omitempty"`
		} `yaml:"prefetch"`
		Consume struct {
			AutoAck   bool `yaml:"autoack, omitempty"`
			Exclusive bool `yaml:"exclusive, omitempty"`
			NoLocal   bool `yaml:"nolocal, omitempty"`
			NoWait    bool `yaml:"nowait, omitempty"`
		} `yaml:"consume"`
	} `yaml:"client"`
	Templates struct{
		Default string `yaml:"default, omitempty"`
		Modules map[string]string `yaml:"modules, omitempty"`
	} `yaml:"templates, omitempty"`
}

// Check config file and build configuration array.
func ConfigSetup(consumer string, configDir string) Config {

	// Set default config values.
	defaultConfigSet()
	viper.SetDefault("consumerTag", consumer)

	// YML based configuration.
	viper.SetConfigType("yaml")

	// Set file name.
	viper.SetConfigName(consumer)

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
			InitErrorHandler("Config file does not exist", err)
		} else {
			// Config file was found but another error was produced
			InitErrorHandler("Fatal error config file:", err)
		}
	}

	// Specific case for mandrill queue type, set default value
	// for template which fallback to blank.
	if viper.GetString("type") == "mandrill" {
		viper.SetDefault("templates.default", "blank")
	}

	// Define configuration variable.
	var config Config

	err := viper.Unmarshal(&config)
	InitErrorHandler("unable to decode into config struct", err)

	// Check supported queue type.
	if config.Type != "mandrill" && config.Type != "mailchimp" {
		InputErrorHandler("you are using unsupported provider type")
	}

	// Channel is required.
	if config.QueueName == "" {
		InputErrorHandler("queue name which you want to consume is required to be defined in yaml file. Yaml key: queue_name")
	}

	return config
}

// Set default values.
func defaultConfigSet() {

	// We set default here, maybe some providers does not require API key.
	viper.SetDefault("apiKey", "")

	// Set rabbitmq settings.
	viper.SetDefault("client.uri", "amqp://guest:guest@localhost:5672")

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
