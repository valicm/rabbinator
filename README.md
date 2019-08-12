# rabbinator
My first lines of code in GO language. Feedback, pull requests appreciated!

# About 
`rabbinator` is a RabbitMQ consumer (worker) written in GO with support for processing Mailchimp and Mandrill 
queue items.

It works out of the box with queue items sent from Drupal based RabbitMQ module (publisher).
But if you publish messages to RabbitMQ in the same format, you can use it.

This implementation does not cover entire Mailchimp / Mandrill API, but rather
this package support only specific operations.
* Mandrill - sending email trough `template send` API call (https://mandrillapp.com/api/docs/messages.JSON.html#method=send-template)
* Mailchimp - subscribe existing or adding new member trough `Add or update a list member` endpoint (https://developer.mailchimp.com/documentation/mailchimp/reference/lists/members/#edit-put_lists_list_id_members_subscriber_hash)

RabbitMQ connection and channel information are configurable trough YAML file.

# Requirements
* RabbitMQ instance with already set channels (this package only consume, does not publish messages)
* Queue items sent to RabbitMQ in a specific format (see details below)
* Mandrill or/and Mailchimp API key
* Configuration provided in YAML format

# Usage
Basic usage example: `./rabbinator -consumer=my_consumer_config_file`

* required flag `consumer` - used to load configuration YAML file, and to set consumer tag in RabbitMQ
* configuration file (`my_consumer_config_file.yml`) is inside a current directory or `$HOME/.rabbinator`
* optional flag `config` - if you want to specify a different location for configuration file then is provided by default
* minimum configuration file needed to run `rabbinator` and consume messages:
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
In Mandrill case, you can specify multiple templates inside the YAML config file. 

With that setup, you can set multiple consumers per channel if you want,
using different names for configuration files.

# Features
* sending mail through Mandrill
* subscribe new and update existing member on Mailchimp
* RabbitMQ acknowledgment based on response and retries in case of failure
* Syslog logs for all failed actions for easier debugging
    
# Configuration
See examples under `examples/config` folder.

* Required yaml keys: `type`, `apiKey`, `queueName`
* Option to adjust RabbitMQ connection - see `examples/config/full_config_example.yml`
* Add Mandrill template mapping - see `examples/config/mandrill_templates_example.yml`

# Drupal project notes
* Supported Drupal 8 queue items sent trough RabbitMQ Drupal module.
* It should work with queue items sent from Drupal 7
* Mailchimp should work with Drupal community module version (I am running slightly custom integration)
* Mandrill works with queue items sent from community module (https://www.drupal.org/project/mandrill)

### Drupal specifics setup
* Disable consuming RabbitMQ items on Drupal cron. Even If you don't use this binary, you should still disable it.
Bootstraping entire Drupal to consume queue is exact opposite of the reason why in first place you are using RabbitMQ.
```
/**
 * Implements hook_queue_info_alter().
 */
function lfi_performance_queue_info_alter(&$queues) {

  // Disable mandrill queue workers. We are processing consumer outside Drupal.
  if (isset($queues['mandrill_queue'])) {
    $queues['mandrill_queue']['cron']['time'] = 0;
  }
}

```

* Mandrill - `google_analytics_campaign` is a string in Drupal implementation. Mandrill declared as string|array. 
The library which this package utilize declared that field as map (array). 
You could trough hook_mandril_mail_alter transform it to an array. (it will work anyway, but throwing a notice)

# RabbitMQ message format
This package is primarily built upon queue items for Mandrill and Mailchimp which are published to RabbitMQ
from Drupal 8 instance using `https://www.drupal.org/project/rabbitmq` module.

You can check examples of how data are structured on sending:

Mandrill `https://git.drupalcode.org/project/mandrill/blob/8.x-1.x/src/Plugin/Mail/MandrillMail.php#L235`

Mailchimp `https://git.drupalcode.org/project/mailchimp/blob/8.x-1.x/mailchimp.module#L386`

_(note that $interests should always be object, avoid using an empty array - it is wrong)_

# Credits
* Using `https://github.com/spf13/viper` for configuration management
* Using `github.com/bkway/gochimp` for communication with Mandrill / Mailchimp API
