# Service Check
The service check is used to perform various checks against a service on an operating system. Until this functionality is brought to parity, the checks supported are different between Linux and Windows.

Linux supports different Service Managers (`systemd`, `init`, etc). Currently `systemd` is the only supported Service Manager and is the default.

## Service Installed and Running Check
Checking that a named service is installed and running is the only check supported by both Linux and Windows.

### Linux
This check only verifies the service is running. Note that `--name (-n)` is required and the `--manager` defaults to `systemd`.
```
check_service --name sshd
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

The functionality depends on the command line flags used and can be easily inferred based on the flags present.
* `--name (-n)` : The service name. Required.
* `--state (-s)` : Validate the service is in the named state
* `--user (-u)` : Validate the service is started by the named user.
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