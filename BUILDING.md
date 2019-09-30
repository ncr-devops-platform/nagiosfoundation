# Building nagiosfoundation
In addition to __Go__ and __Make__, building requires a couple of utilities.
* [__gödel__](https://github.com/palantir/godel) build tool
* [__dep__](https://github.com/golang/dep) dependency management tool

During the build gödel will automatically install itself at `$HOME/.godel` but dep must be installed beforehand using the [dep installation instructions](https://github.com/golang/dep#installation) or follow the steps below for Linux.

### Create workspace
```
cd $GOPATH
mkdir -p src/github.com/ncr-devops-platform/
cd src/github.com/ncr-devops-platform/
```

### Clone project
```
git clone https://github.com/ncr-devops-platform/nagios-foundation.git nagiosfoundation
```

### Set project directory
```
cd nagiosfoundation
```

### Install dep
For Linux. Use the [dep installation instructions](https://github.com/golang/dep#installation) for other platforms.
```
mkdir $GOPATH/bin
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

### Update dependencies
```
$GOPATH/bin/dep ensure -v
```

### Build
```
make
```

### Artifacts
The build artifacts can be found in `out/build`.

## Releases
The `Makefile` uses [`ghr`](https://github.com/tcnksm/ghr) to create the GitHub Release and upload artifacts. This would normally be done by the project owner and requires a GitHub API token. Refer to the [ghr project](https://github.com/tcnksm/ghr) for more information and install `ghr` with
```
go get -u github.com/tcnksm/ghr
```
The Makefile target is `release` and is used when lauching make with
```
make release
```

## Individual Commands
When developing it may be desirable to build a single command rather than the entire project. There are a few ways of doing this. Choose the one appropriate for you.

### Gödel
The `Makefile` uses Gödel for building so this method is much like the `Makefile` method.

To build a command for all `os-archs` as listed in the `godel/config/dist-plugin.yml` configuration.
```
./godelw build check_process
```

### Go Build
The native `go build` command can also be used to build a command for a single operating system (default is your native OS). Note this command will drop the executable into the current directory. To change this, execute `go help build` and look for the `-o` option.
```
go build github.com/ncr-devops-platform/nagiosfoundation/cmd/check_process
```

To build for a different OS, use the `GOOS` and `GOARCH` environment variables. For example, to build an executable for Windows on Linux, use:
```
GOOS=windows GOARCH=amd64 go build -o check_process.exe github.com/ncr-devops-platform/nagiosfoundation/cmd/check_process
```

### Adding a New Check

This project uses a `Makefile` for builds. In turn, the `Makefile` uses [Gödel](https://github.com/palantir/godel) as the next step in the chain. Adding a new check consists of three main steps. Keep in mind these are general guidelines and don't necessarily need to be done in order.

1. Add the main logic for your check to `lib/app/nagiosfoundation`. Consider this directory as containing the APIs. As APIs, any code placed here should: gracefully handle all errors, not exit (such as calling `os.Exit()`), and not send output to the terminal. Input into the API will generally originate from command-line options and parameters from the user interface and passed as parameters into the API. Output from the API consists of a result check text string using `resultMesssage()` or one of the functions in the `nagiosformatters` package as well as a numeric representation of the check result which should be one of `statusTextOK`, `statusTextWarning`, or `statusTextCritical` which are `const`s found in `lib/app/nagiosfoundation/resultmessage.go`. Bonus points for implementing unit tests for new API code. Tests for the API code are especially helpful for us maintainers as it helps the CI process verify defects haven't been introduced. Even just exercising the code is helpful.

2. Add user facing logic to `cmd/<check_name>`. Directories in `cmd/` contain the user interfaces which in this case consist of a command-line interface implemented by the [Cobra framework](https://github.com/spf13/cobra). User input consists of command-line options and parameters, with command output being the check result text received from the API and sent to the terminal, and the error code received from the API and returned to the shell with `os.Exit()`. The new check should output version information on demand. Make sure this is enabled with a call to `initcmd.AddVersionCommand()` because the Travis builds check this with `scripts/validate-version.sh` and will fail the build if a proper version is not output. Unit tests for this code is desireable but not as important as for API code.

3. Add the new check to the build process by editing `godel/config/dist-plugin.yaml` and adding your new check. Generally this is as simple as copying the entry for an existing check and making the obvious changes to the check names and perhaps build targets.