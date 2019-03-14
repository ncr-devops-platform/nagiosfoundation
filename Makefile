# godel_version is only used to ensure all the godel bits
# have been downloaded before determining the project-version.
godel_version := $(shell ./godelw version)
version := $(shell ./godelw project-version)
package_path = ./out/package
package_version = $(package_path)/$(version)
platforms = windows-amd64 linux-amd64 windows-386 linux-386

package: $(platforms)

clean:
	./godelw clean
	rm -rf $(package_path)
	rm -rf coverage.txt
	rm -rf coverage.html

test:
	go test -v ./...

test-godel:
	./godelw test

coverage:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

release: clean package
	ghr $(version) $(package_path)

build:
	./godelw build

$(platforms): build
	$(eval package_bin = $(package_version)/$@/bin)
	mkdir -p $(package_bin)
	ln ./out/build/*/$(version)*/$@/* $(package_bin)/.
	tar -zcvf $(package_path)/nagiosfoundation-$@-$(version).tgz -C $(package_version)/$@ bin
	(cd $(package_path) && sha512sum nagiosfoundation-$@-$(version).tgz) > $(package_path)/nagiosfoundation-$@-$(version)-sha512.txt
	rm -rf $(package_version)
