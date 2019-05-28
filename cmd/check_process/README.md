# Process Check
The process check attempts to find a process by name specified with the `--name (-n) flag`. The result of the check depends on the value of the `--type (-t)` flag. If the `--type` flag is not specified, the default is `running`. Valid types are:
* `running`: If the process is found, the check returns an `OK` result, otherwise it returns `CRITICAL`.
* `notrunning` If the flag is not found, the check returns `OK` result, otherwise it returns `CRITICAL`.

## Process Running
```
check_process --name bash --type running
```

## Process Not Running
```
check_process --name invalidname --type notrunning
```
