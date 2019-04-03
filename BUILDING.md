# Building nagiosfoundation
In addition to __Go__ and __Make__, building requires the [__gödel__](https://github.com/palantir/godel) build tool which will automatically install itself during the build..

The build also uses [Go Modules](https://blog.golang.org/using-go-modules) for which Go 1.11 and Go 1.12 include preliminary support. In Go 1.13 this will be default for all Go development.

### Clone project
Clone the project somewhere outside of your [Go Workspace](https://golang.org/doc/code.html#Workspaces) (a requirement of Go Modules).

```
git clone https://github.com/ncr-devops-platform/nagios-foundation.git
```

### Set project directory
```
cd nagios-foundation
```

### Build
```
make
```

### Artifacts
The build artifacts can be found in `out/build` and full build assets are in `out/package`.

## Releases
Releases are normally performed through Travis CI by tagging a commit. As a secondary method, the `Makefile` can use [`ghr`](https://github.com/tcnksm/ghr) to create the GitHub Release and upload artifacts. This would normally be done by the project owner and requires a GitHub API token. Refer to the [ghr project](https://github.com/tcnksm/ghr) for more information and install `ghr` with
```
go get -u github.com/tcnksm/ghr
```
The Makefile target is `release` and is used when lauching make with
```
make release
```

## Individual Commands
When developing it may be desirable to build a single command rather than the entire project. There are a few ways of doing this. Choose the one appropriate for you.

### Go Module Configuration
Because this project uses Go Modules for dependencies, before performing individual builds execute the following to trigger proper functionality if this project was placed in a Go Workspace, or run it anyway if you're not sure what this means.
```
export GO111MODULE=1
```

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

