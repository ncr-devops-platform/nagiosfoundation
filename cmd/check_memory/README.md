# Memory Check
The memory check (`check_memory`) checks the available memory as reported by the OS. It queries the OS for the amount of available memory, the amount of free memory, then calculates a memory used percentage. That memory used percentage is then compared against the `--warning` and `--critical` thresholds and an appropriate check result is output.

## Flags
* `--warning`: The percentage of used memory required to trigger a warning condition. Default `85`.
* `--critical`: The percentage of used memory required to trigger a critical condition. Default `95`.
* `--metric_name`: The name used in the nagios portion of the message output. Default `memory_used_percentage`.

## Examples
Issue a warning if memory usage is over 50% and critical if usage is over the default of 95%.
```
check_memory --warning 50
```