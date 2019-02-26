# Building nagiosfoundation
In addition to __Go__ and __Make__, building requires a couple of utilities.
* [__gödel__](https://github.com/palantir/godel) build tool
* [__dep__](https://github.com/golang/dep) dependency management tool

During the build gödel will automatically install itself at `$HOME/.godel` but dep must be installed beforehand using the [dep installation instructions](https://github.com/golang/dep#installation) or follow the steps below for Linux.

### Create workspace
```
cd $GOPATH
mkdir -p nagios-foundation/src/github.com/jkerry/
cd nagios-foundation/src/github.com/jkerry/
```

### Clone project
```
git clone https://github.com/jkerry/nagios-foundation.git nagiosfoundation
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
