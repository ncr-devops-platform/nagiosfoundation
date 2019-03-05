# Performance Counter Check
The performance counter check is Windows only. It retreives a Windows Performance Counter (`-counter_name`) and compares it to `-critical` and `-warning` then outputs an appropriate responsed based on the check. Many flags make this check quite configurable.

The defaults for this check have the `-critical` and `-warning` values set to `0`, and the counter value retrieved is compared to be lesser than those values. Generally a counter value will be `> 0`, causing this check to generally emit an `OK` response when using these defaults.

* `-counter_name`: The Performance Counter to fetch. No default and therefore it must specified.
* `-greater_than`: If set, the Performance Counter value is compared to be greater than the `-critical` and `-warning` values. The default value of `false` causes the comparison to be less than the `-critical` and `-warning` values. Do not specify a flag argument when using.
* `-critical`: The value the performance counter is compared with to determine a `CRITICAL` response. The default is `0`.
* `-warning`: The value the performance counter is compared with to determine a `WARNING` response. The default is `0`.
* `-polling_attempts`: The number of times to attempt retrieval of the Performance Counter. Default is `2`.
* `-polling_delay`: The delay, in seconds, between polling attempts. Default is `1`.
* `-metric_name`: String output in the response message in a Nagios format and is suitable for machine parsing. There is no default.

## Examples
Specify only the required counter name
```
check_performance_counter -counter_name "\IPv4\Datagrams/sec"
```

## Helpful Hints
To get a list of Performance Counters available on a system use
```
TypePerf.exe –q
```
