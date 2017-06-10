
# do not specify a full path for go since travis will fail
GO = GOGC=off go
GOFLAGS = -ldflags "-X main.version=$(shell git describe --tags)"
GOVENDOR = $(shell which govendor)

all: build test

help:
	@echo "build     - go build"
	@echo "install   - go install"
	@echo "test      - go test"
	@echo "gofmt     - go fmt"
	@echo "linux     - go build linux/amd64"
	@echo "release   - build/release.sh"
	@echo "homebrew  - build/homebrew.sh"
	@echo "buildpkg  - build/build.sh"
	@echo "pkg       - build, test and create pkg/olb.tar.gz"
	@echo "clean     - remove temp files"

build: checkdeps
	$(GO) build -i $(GOFLAGS)
	$(GO) test -i ./...

test: checkdeps
	$(GO) test -v -test.timeout 15s `go list ./... | grep -v '/vendor/'`
	@if [ $$? -eq 0 ] ; then \
	    echo "All tests PASSED" ; \
    else \
	    echo "Tests FAILED" ; \
	fi

checkdeps:
	[ -x "$(GOVENDOR)" ] || $(GO) get -u github.com/kardianos/govendor
	govendor list +e | grep '^ e ' && { echo "Found missing packages. Please run 'govendor add +e'"; exit 1; } || : echo

gofmt:
	gofmt -w `find . -type f -name '*.go' | grep -v vendor`

linux:
	GOOS=linux GOARCH=amd64 $(GO) build -i -tags netgo $(GOFLAGS)

install:
	$(GO) install $(GOFLAGS)

pkg: build test
	rm -rf pkg
	mkdir pkg
	tar czf pkg/olb.tar.gz olb

release: test
	build/release.sh

homebrew:
	build/homebrew.sh

codeship:
	go version
	go env
	unzip -o -d ~/bin ~/vault.zip
	vault --version
	cd ~/src/github.com/millisecond/olb && make test

olb-builder: make-olb-builder push-olb-builder

make-olb-builder:
	docker build -t olblb/olb-builder --squash build/olb-builder

push-olb-builder:
	docker push olblb/olb-builder

buildpkg: test
	build/build.sh

clean:
	$(GO) clean
	rm -rf pkg

.PHONY: build linux gofmt install release docker test homebrew buildpkg pkg clean
