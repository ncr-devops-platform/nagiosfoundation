# Uptime Check
The uptime check (`check_uptime`) system uptime. The uptime value is then compared against the `--warning` and `--critical` thresholds and an appropriate check result is output.

## Flags
* `--warning`: The desired time (in seconds(s), minutes(m), or hours(h)) of uptime required to trigger a warning condition. Default is 72h.
* `--critical`: The desired time (in seconds(s), minutes(m), or hours(h)) of uptime required to trigger a critical condition. Default  is 1 week (168h).
* `--metric_name`: The name used in the Nagios portion of the message output. Default `current_system_uptime`.

## Examples
Issue a warning if uptime is over 72 hours and critical if uptime is over the default of 1 week.
```
check_uptime --warning 72h --critical 168h
```