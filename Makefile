SHELL:=/bin/bash

up:
	go get -u all
	go mod tidy
tag:
	{ \
		set -e; \
    	TAG=v$(shell cat ./version); \
    	ORIGIN="origin"; \
    	git tag -d "$$TAG" || true; \
        git push -d "$$ORIGIN" "$$TAG" || true; \
        git tag "$$TAG"; \
        git push --tags "$$ORIGIN" "$$TAG"; \
	};