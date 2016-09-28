
all: install

install: glide-install
	go install

build:
	go build

glide-update:
	glide update

glide-install:
	glide install
