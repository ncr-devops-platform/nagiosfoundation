# Process memory check
The process memory check (`check_process_memory`) checks the amount of memory consumed by a process as reported by OS. Retrieved memory percentage is compared against `--warning` and `--critical` thresholds and an appropriate check result is output.

## Linux process memory check
On linux machines `ps` command will be used to retrieve memory percentage. Valid PIDs will be discovered based on `/proc/[pid]/stat` files.

## Windows process memory check
On windows machines Windows Management Instrumentation (WMI) will be used. WMI will be queried for each running process memory consumption. Processes matching `<PROCESS_NAME>.exe` **and children** of those processes will be added to the calculation.

> NOTE: It is possible that parent id of the process refers to a different (reused) PID, which theoretically may yield incorrect results.

> This check uses WorkingSetSize field of a process to calculate percentage. Sum of all processes' WorkingSetSize should be equal to what is observed in Task Manager as total used memory. Sum of a particular processes' WorkingSetSize MIGHT NOT be equal to what is reported as total consumed memory of a process in Task Manager.

## Flags
> Run this check with help command to get latest information, i.e. `check_process_memory help`
* `--warning`: The percentage of used memory required to trigger a warning condition. Default `85`.
* `--critical`: The percentage of used memory required to trigger a critical condition. Default `95`.
* `--metric_name`: The name used in the nagios portion of the message output. Default `memory_used_process_percentage`
* `--process_name`: **Required.** The name of the process for which to query memory usage. Process name is **case-sensitive**.

## Examples
Issue a warning if memory usage is over 50% for chrome
```
check_process_memory --process_name chrome --warning 50
```