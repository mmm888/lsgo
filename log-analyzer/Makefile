NAME := log-analyzer
APP := github/mmm888/mycmd/log-analyzer

.PHONY: build
build: vendor
	go build -o bin/$(NAME)

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: vendor
vendor: dep-install
	dep ensure -vendor-only=true

.PHONY: dep-install
dep-install:
ifeq ($(shell type dep 2> /dev/null),)
	go get -u github.com/golang/dep/...
endif

.PHONY: dep-ensure
dep-ensure: dep-install
	dep ensure

test: 
	go test $(APP)/... -cover
