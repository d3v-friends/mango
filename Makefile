SHELL:=/bin/bash

upmod:
	go get -u all
	go mod tidy
tag:
	sh script/tag.sh $(shell cat ./version)