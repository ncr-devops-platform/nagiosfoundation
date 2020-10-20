# Performance Counter Check
The performance counter check is Windows only. It retreives a Windows Performance Counter (`--counter_name`) and compares it to `--critical` and `--warning` then outputs an appropriate responsed based on the check. Many flags make this check quite configurable.

The defaults for this check have the `--critical` and `--warning` values set to `0`, and the counter value retrieved is compared to be lesser than those values. Generally a counter value will be `> 0`, causing this check to generally emit an `OK` response when using these defaults.

* `--counter_name (-n)`: The Performance Counter to fetch. No default and therefore it must be specified. This can take multiple counters as an input and generates the data for the given counters. Multiple counters can be passed seperated by a comma as shown in the example below.
When passing multiple counters, if the flags like `-m or -g or -w or -c` are passed to the perf counter then it throws an error. Therefore, multiple performance counter is used for data retrival but it doesn't support the comparison of warning and critical values just like single performance counter.

* `--greater_than (-g)`: If set, the Performance Counter value is compared to be greater than the `--critical` and `--warning (-w)` values. The default value of `false` causes the comparison to be less than the `--critical` and `--warning` values. Do not specify a flag argument when using.
* `--critical (-c)`: The value the performance counter is compared with to determine a `CRITICAL` response. The default is `0`.
* `--warning (-w)`: The value the performance counter is compared with to determine a `WARNING` response. The default is `0`.
* `--polling_attempts (-a)`: The number of times to attempt retrieval of the Performance Counter. Default is `2`.
* `--polling_delay (-d)`: The delay, in seconds, between polling attempts. Default is `1`.
* `--metric_name (-m)`: String output in the response message in a Nagios format and is suitable for machine parsing. There is no default.

## Examples
Specify only the required counter name
```
check_performance_counter --counter_name "\IPv4\Datagrams/sec"
```

## Example for passing multiple inputs for performance counter

Pass all the commands for which you want to get the data
```
check_performance_counter -n "\LogicalDisk(C:)\% Free Space, \LogicalDisk(C:)\FreeMegaBytes, \system\system up time, \TCPv4\ConnectionsEstablished"
```
This will print us the data(values) for the above commands in the order.
If all the counters defined are available, it prints the values and then the overall check status will be a success(0)
If any of the counters defined are not available, it displays the error and overall check status will be critical(2)

## Helpful Hints
To get a list of Performance Counters available on a system use
```
TypePerf.exe –q
```
