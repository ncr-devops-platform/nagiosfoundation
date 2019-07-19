# Uptime Check
The uptime check (`check_uptime`) system uptime. The uptime value is then compared against the `--warning` and `--critical` thresholds and an appropriate check result is output.

## Flags
* `--warning`: The desired time (in seconds) of uptime required to trigger a warning condition. Default is 72 hours (`259200`).
* `--critical`: The desired time (in seconds) of uptime required to trigger a critical condition. Default  is 1 week (`604800`).
* `--metric_name`: The name used in the Nagios portion of the message output. Default `current_system_uptime`.

## Examples
Issue a warning if uptime usage is over 50% and critical if usage is over the default of 95%.
```
check_uptime --warning 7200 --critical 14200
```