# Service Check
The service check is used to perform various checks against a service on an operating system. Until this functionality is brought to parity, the checks supported are different between Linux and Windows.

Because Linux supports different Service Managers (`systemd`, `init`, etc), a service manager must be specified. Currently `systemd` is the only supported Service Manager and it must be specified with the `-manager` flag.

## Service Installed and Running Check
This is the only check supported by both Linux and Windows.

### Linux
This check only verifies the service is running. Note that both `-name` and `-manager` are required.
```
check_service -name sshd -manager systemd
```

### Windows
Because the various states are supported, the running state must be specified with the `-state` flag.
```
check_service.exe -name audiosrv -state running
```

## Windows Extended Features
The Windows version of this check supports several additional features.
* A service exists, outputting the current service state and user.
* A service exists and is in a specified state.
* A service exists and is started by a specified user.
* A service exists, is in a specified state, and is started by a specified user.

The functionality depends on the command line flags used and can be easily inferred based on the flags present.
* `-name`: The service name. Required.
* `-state`: Validate the service is in the named state
* `-user`: Validate the service is started by the named user.

### Service Exists
```
check_service.exe -name audiosrv
```
### Service Exists and in State
```
check_service.exe -name audiosrv -state stopped
```
### Service Exists and Started by User
```
check_service.exe -name audiosrv -user "NT AUTHORITY\LocalService"
```
### Service Exists, in State, and Started by User
```
check_service.exe -name audiosrv -state running -user "NT AUTHORITY\LocalService"
```