# Service Check
The service check is used to perform various checks against a service on an operating system. Until this functionality is brought to parity, the checks supported are different between Linux and Windows.

Linux supports different Service Managers (`systemd`, `init`, etc). Currently `systemd` is the only supported Service Manager and is the default.

## Common Checks
Both Linux and Windows support checking that a named service is running and output of the current state in a nagios format.

### Linux
This check supports
* Verify a service is running
* Return the state of the service as a nagios formatted result

The functionality depends on the command line flags used and can be easily inferred based on the flags present.
* `--name (-n)` : The service name. Required.
* `--current-state (-c)` : Output the service state in nagios output

### Service Running
To verify a service is running, use the `--name (-n)` option. Note that `--name (-n)` is required and the `--manager` defaults to `systemd`.
```
check_service --name sshd
```

### Return the State of a Service
Use the `--current_state (-c)` along with the `--name (-n)` option to return the state of a service. Linux supports two states:
* 0 - Not running
* 1 - Running

Some examples are

**Service Running**

```
$ check_service --name sshd --current_state
CheckService OK - sshd in a running state | service_state=1 service_name=sshd
```

**Service Not Running**

```
$ check_service --name sshd --current_state
CheckService CRITICAL - sshd not in a running state (State: inactive) | service_state=0 service_name=sshd
```

### Windows
Various states are supported and the `--state (-s)` flag is used to specify them. The default is to check that the service exists.
```
check_service.exe --name audiosrv --state running
```

## Windows Extended Features
The Windows version of this check supports several additional features.
* A service exists, outputting the current service state and user.
* A service exists and is in a specified state.
* A service exists and is started by a specified user.
* A service exists, is in a specified state, and is started by a specified user.
* Returning the state of a service as a nagios formatted result

The functionality depends on the command line flags used and can be easily inferred based on the flags present.
* `--name (-n)` : The service name. Required.
* `--state (-s)` : Validate the service is in the named state
* `--user (-u)` : Validate the service is started by the named user.
* `--current-state (-c)` : Output the Windows service state in nagios output
* `--manager (-m)` : Specify a service manager. `wmi` and `svcmgr` are supported. The default is `wmi`.

## Windows Service Manager
The Windows version of this check supports two methods of retrieving service data.

**`wmi`**: Uses [Windows Management Instrumentation](https://docs.microsoft.com/en-us/windows/desktop/wmisdk/wmi-start-page) to retrieve service data. This method does not require any special user privileges and is the default.

**`svcmgr`**: Uses the [Windows Control Manager](https://docs.microsoft.com/en-us/windows/desktop/services/service-control-manager) to retrieve service data. This method requires sufficient user privileges to access the control manager.

### Service Exists
```
check_service.exe --name audiosrv
```
### Service Exists and in State
```
check_service.exe --name audiosrv --state stopped
```
### Service Exists and Started by User
```
check_service.exe --name audiosrv --user "NT AUTHORITY\LocalService"
```
### Service Exists, in State, and Started by User
```
check_service.exe --name audiosrv --state running --user "NT AUTHORITY\LocalService"
```

### Return the State of a Service
```
./check_service.exe --name audiosrv --current_state
CheckService OK - audiosrv is in a Running state | service_state=0 service_name=audiosrv
```

### Return the State of Non-existing Service
```
./check_service.exe --name fakeservice --current_state
CheckService OK - fakeservice does not exist | service_state=255 service_name=fakeservice
```

## Service State Numbers
The state numbers returned using the `--current_state (-c)` option are:
* 0 - Running
* 1 - Paused
* 2 - Start Pending
* 3 - Pause Pending
* 4 - Continue Pending
* 5 - Stop pending
* 6 - Stopped
* 7 - Unknown
* 255 - No such service