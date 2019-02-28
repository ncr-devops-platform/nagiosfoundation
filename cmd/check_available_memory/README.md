# Available Memory Check
The available memory check (`check_available_memory`) checks the available memory as reported by the OS.

## Intended Purpose
The probable intended purpose of the available memory check is to query the OS for the amount of available memory, the amount of free memory, then calculate a memory used percentage. That memory used percentage is then compared against the `-warning` and `-critical` thresholds and an appropriate check result is output.

### Flags
* `-warning`: The percentage of used memory required to trigger a warning condition.
* `-critical`: The percentage of used memory required to trigger a critical condition.
* `-metric_name`: The name used in the nagios portion of the message output.

## Implemented Purpose
Someone looking through the code may notice what is currently implemented doesn't match what is probably the intended purpose. The details here are to help someone contributing to understand.

The available memory check currently queries the OS for the amount of available memory (megabytes in Windows, kilobytes on Linux) and compares that against the `-warning` and `-critical` thresholds. These thresholds are `85` and `95` respectively. Because the amount of free memory available on an active OS is generally more than those thresholds, the check usually returns on OK status.

For Linux, the check scrapes the `free` command and pulls the swap total instead of memory free.