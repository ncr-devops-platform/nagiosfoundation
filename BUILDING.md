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

