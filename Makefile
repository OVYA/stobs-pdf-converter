.PHONY: all
all: install

.PHONY: install
install: glide-install
	go install

.PHONY: build
build: vendor
	go build

.PHONY: glide-update
glide-update:
	glide update

.PHONY: glide-install
glide-install: vendor
	glide install

vendor:
	glide install

.PHONY: package
package:
	go build
	mv stobs package/stobs/usr/bin/
	dpkg-deb --build package/stobs
