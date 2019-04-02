# nagios-foundation

[![Build Status](https://travis-ci.org/ncr-devops-platform/nagios-foundation.svg?branch=master)](https://travis-ci.org/ncr-devops-platform/nagios-foundation)
[![codecov](https://codecov.io/gh/ncr-devops-platform/nagios-foundation/branch/master/graph/badge.svg)](https://codecov.io/gh/ncr-devops-platform/nagios-foundation)

A suite of Nagios style checks and metrics covering the basic needs for monitoring in a Sensu-like system.


## List of Checks
* [CPU](https://github.com/ncr-devops-platform/nagios-foundation/blob/master/cmd/check_cpu/README.md)
* [Memory](https://github.com/ncr-devops-platform/nagios-foundation/blob/master/cmd/check_available_memory/README.md)
* [Performance Counter](https://github.com/ncr-devops-platform/nagios-foundation/blob/master/cmd/check_performance_counter/README.md)
* [Process](https://github.com/ncr-devops-platform/nagios-foundation/blob/master/cmd/check_process/README.md)
* [Service](https://github.com/ncr-devops-platform/nagios-foundation/blob/master/cmd/check_service/README.md)
* [User and Group](https://github.com/ncr-devops-platform/nagios-foundation/blob/master/cmd/check_user_group/README.md)

## Using
Use this collection of applications as [Sensu Go Checks](https://docs.sensu.io/sensu-go/5.5/reference/checks/) in your Sensu deployment. For example, to check every 60 seconds that the signage application is running on a remote kiosk where the Sensu Agent is subscribed to `signage`, run:

```
$ cat << EOF | sensuctl create
{
  "type": "CheckConfig",
  "api_version": "core/v2",
  "metadata": {
    "name": "process_signage",
    "namespace": "default"
  },
  "spec": {
    "command": "check_process --name signage_app",
    "interval": 60,
    "publish": true,
    "subscriptions": [
      "signage"
    ]
  }
}
EOF
```

---

## Building and Contributing
See [Build Instructions](https://github.com/ncr-devops-platform/nagios-foundation/blob/master/BUILDING.md)

## Testing and Code Coverage
See [Testing](https://github.com/ncr-devops-platform/nagios-foundation/blob/master/TESTING.md)
