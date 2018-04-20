GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOPACKAGES=$(shell go list ./... | grep -v /vendor/)

fmt:
	@if [ -n "$$(gofmt -l ${GOFILES})" ]; then echo 'Please run gofmt -l -w on your code.' && exit 1; fi

copyright:
	@${GOPATH}/src/github.com/IBM/portieris/scripts/copyright.sh

test:
	$(GOPATH)/bin/gotestcover -v -coverprofile=cover.out ${GOPACKAGES}

