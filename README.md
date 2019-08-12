# rabbinator
My first lines of code in GO language. Don't judge :-D

# About 
`rabbinator` is a RabbitMQ consumer (worker) written in GO with support for processing Mailchimp and Mandrill 
queue items.

It works out of box with queue items sent from Drupal based RabbitMQ module (publisher).
But if you publish messages to RabbitMQ in same format you can use it.

This implementation does not cover entire Mailchimp / Mandrill API, but rather
this package support only specific operations.
* Mandrill - sending email trough template send API call (https://mandrillapp.com/api/docs/messages.JSON.html#method=send-template)
* Mailchimp - subscribe existing or adding new member trough `Add or update a list member
` endpoint (https://developer.mailchimp.com/documentation/mailchimp/reference/lists/members/#edit-put_lists_list_id_members_subscriber_hash)

RabbitMQ connection and channel information are configurable trough YAML file.

# Requirements
* RabbitMQ instance with already set channels (this package only consume, does not publish messages)
* Queue items sent to RabbitMQ in specific format (see details below)
* Mandrill or / and Mailchimp API key
* Configuration provided in yaml format

# Usage
Basic usage example: `./rabbinator -consumer=my_consumer_config_file`

* reguired flag `consumer` - used to load configuration yml file,and to set consumer tag in RabbitMQ
* configuration file (`my_consumer_config_file.yml`) is inside current directory or `$HOME/.rabbinator`
* optional flag `config` - if you want to specify different location for configuration file then is provided by default
* minimum configuration file needed to run rabbinator and consume messages:
```
type: "mandrill"
apiKey: "your_api_key"
queueName: "rabbitmq_name_of_mandrill_queue"
```

It is required to define:
 * type: `mandrill` or `mailchimp`
 * apiKey: your account apikey
 * queueName: RabbitMQ channel which you want to consume
    
Yaml configuration provides entire mapping for RabbitMQ connection & settings.
Also in Mandrill case you can specify multiple templates. 

With that setup, you can set multiple consumers per channel if you want,
using different names for configuration files.
    
# Configuration
See examples under `examples/config` folder.

* Required yaml keys: `type`, `apiKey`, `queueName`
* Option to adjust RabbitMQ connection - see `examples/config/full_config_example.yml`
* Add Mandrill template mapping - see `examples/config/mandrill_templates_example.yml`

# Drupal project notes
* it should work with queue items sent from Drupal 7
* mailchimp should work with Drupal community module version (I am running custom integration)
* mandrill works with queue items sent from community module (https://www.drupal.org/project/mandrill)

# RabbitMQ message format
Primary this package is built upon queue items for Mandrill and Mailchimp which are published to RabbitMQ
from Drupal 8 instance using `https://www.drupal.org/project/rabbitmq` module.

# Credits
* Using `https://github.com/spf13/viper` for configuration management
* Using `github.com/bkway/gochimp` for communication with Mandrill / Mailchimp API