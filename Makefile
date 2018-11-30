version := $(shell git describe --tags)
platforms = windows-amd64 linux-amd64


package: $(platforms)

clean:
						./godelw clean
						rm -rf ./out/package

$(platforms):
						./godelw build
						mkdir -p ./out/package/$(version).dirty/$@/bin
						cp ./out/build/*/$(version).dirty/$@/* ./out/package/$(version).dirty/$@/bin/.
						tar -cvf ./out/package/nagiosfoundation-$@-$(version).dirty.tgz -C ./out/package/$(version).dirty/$@ bin

