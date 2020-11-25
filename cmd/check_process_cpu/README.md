# Process CPU check
The process cpu check (`check_process_cpu`) checks the CPU utilization of a process as reported by OS. Retrieved percentage is compared against `--warning` and `--critical` thresholds and an appropriate check result is output.

## Linux process CPU check
On linux machines there are 2 modes of operation supported:
- Global CPU utilization of a process (global usage on OS)
- Per-core CPU utilization of a process

### Global CPU check
Global CPU check uses `top` command in batch mode to find out CPU utilization by a process. Multiple samples will be taken to determine utilization. Returned percentage is normalized for core count (i.e. range is always `[0; 100]`).
> First sample of the `top` command will be ignored as it reports process' lifetime stats on old machines.

### Per-core CPU check
Per-core CPU check uses `pidstat` command. Multiple samples will be taken. Core utilization will be averaged based on amount of non-0 per-core values reported.

For example:
- sample 0: core 0 - 20.0%
- sample 1: core 0 - 10.0%, core 2 - 30.0%
- sample 2: core 1 - 15.0%

Results in:
- core 0 - 15.0%
- core 1 - 15.0%
- core 2 - 30.0%

30.0% will be reported by the check (as the highest core utilization).

## Windows process CPU check
On windows machines only global CPU utilization of a process can be retrieved. Windows Management Instrumentation (WMI) will be used and multiple samples will be taken to calculate average CPU utilization of a process. Returned percentage is normalized for core count (i.e. range is always `[0; 100]`).

## Flags
> Run this check with help command to get latest information, i.e. `check_process_cpi help`
* `--warning`: The percentage of CPU utilization required to trigger a warning condition. Default: `85`.
* `--critical`: The percentage of CPU utilization required to trigger a critical condition. Default: `95`.
* `--metric_name`: The name used in the nagios portion of the message output. Default `process_cpu_percentage`.
* `--process_name`: **Required.** The name of the process for which to query CPU usage. Process name is **case-sensitive**. On windows use executable name without extension (e.g. `chrome`). On linux use command name (e.g. `cbdaemon`)
* `--core`: When included, check will return highest CPU **Core** utilization instead of global utilization. *Only linux machines support this flag; on windows machines this flag will be ignored.*

## Examples
Issue a warning if global CPU utilization is over 50% for chrome
```
check_process_cpu --process_name chrome --warning 50
```

Issue a critical alert if per-core CPU utilization is over 80% for chrome
```
check_process_cpu --process_name chrome --critical 80 --core
```