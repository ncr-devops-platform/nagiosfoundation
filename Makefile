version := $(shell ./godelw project-version)
package_path = ./out/package
package_version = $(package_path)/$(version)
platforms = windows-amd64 linux-amd64 windows-386 linux-386

package: $(platforms)

clean:
	./godelw clean
	rm -rf $(package_path)

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
	rm -rf $(package_version)
