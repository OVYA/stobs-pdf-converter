.PHONY: all
all: install

.PHONY: install
install: glide-install
	go install

.PHONY: build
build:
	go build

.PHONY: glide-update
glide-update:
	glide update

.PHONY: glide-install
glide-install:
	glide install
