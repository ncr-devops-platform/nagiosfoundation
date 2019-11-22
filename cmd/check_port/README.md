# Port Check
The port check attempts to connect to an address specified with the `--address` or `-a` flag and port specified with the `--port` or `-p` flag.

The check can be further refined by using the `--timeout` or `-t` flag to customize the connection wait time and the `--invert` or `-i` flag to return success when a port is not being listened on. Execute `check_port --help` for more details.

With `--invert` set to the default of `false`, `check_port` will return `OK` (exit code of `0`) if the address and port connection is successful. Otherwise it will return `CRITICAL` (exit code `2`).

The Nagios metric name defaults to `listening_port` but can be changed with the `--metric_name` or `-m` option.

## Listener on Address and Port
```
$ check_port --address www.google.com --port 443
CheckPort OK - port 443 on www.google.com using tcp | listening_port=0
```

## No Listener on Address and Port
```
$ check_port --address www.google.com --port 81 --timeout 2
CheckPort CRITICAL - port 81 on www.google.com using tcp (dial tcp 216.58.192.228:81: i/o timeout) | listening_port=2
```

## Inverted No Listener on Address and Port
```
$ check_port --address www.google.com --port 81 --timeout 2 --invert
CheckPort OK - port 81 on www.google.com using tcp (dial tcp 216.58.192.228:81: i/o timeout) | listening_port=0
```
