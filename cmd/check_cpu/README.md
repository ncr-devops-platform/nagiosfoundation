# CPU Check
The CPU check retrieves the CPU load as a percentage and if it is over `--critical` (default 95%), will output a `CRITICAL` response. If a `CRITICAL` response is not output, it compares it to `--warning` (default 85%) and if over, will output a `WARNING` response. Anything else will output an `OK` response.

The `--critical` and `--warning` flags can be set on the command line as desired.

A `--metric_name` flag can also be specified. This string is output in the response message in a Nagios format and is suitable for machine parsing. The default is `pct_processor_time`.

## Examples
Default check values
```
check_cpu
```

Warning level set to 70%
```
check_cpu --warning 70
```

Override the metric name
```
check_cpu --metric_name cpu_percentage
```
