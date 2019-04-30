# Makefile based on https://github.com/Masterminds/glide/blob/master/Makefile

GO                 ?= go
DIST_DIRS          := find * -type d -exec
VERSION            ?= $(shell git describe --tags)
SOURCES            ?= *.go */*.go

# ------------------------------------------------------------------------------

.PHONY: build install clean bootstrap-dist build-all dist

##
## Available make targets
##

# default: show target list
all: help

# ------------------------------------------------------------------------------
## Sources

## Format go sources
fmt:
	$(GO) fmt ./...

## Run vet
vet:
	$(GO) vet ./...

## Run linters
lint:
	golint ./...
	golangci-lint run ./...

## Run tests and fill coverage.out
cov: coverage.out

# internal target
coverage.out: $(SOURCES)
	GIN_MODE=release $(GO) test -test.v -test.race -coverprofile=$@ -covermode=atomic -tags test ./...

## Open coverage report in browser
cov-html: cov
	$(GO) tool cover -html=coverage.out

## Clean coverage report
cov-clean:
	rm -f coverage.*

build:
	${GO} build -o apisite -ldflags "-X main.version=${VERSION}"

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./apisite ${DESTDIR}/usr/local/bin/apisite

clean:
	rm -f ./apisite
	rm -rf ./dist

#test:
#	go test -c -coverpkg=. -tags test

# ------------------------------------------------------------------------------
## Deploy

## Install gox
bootstrap-dist:
	${GO} get -u github.com/Masterminds/gox

## Build dist binaries
build-all:
	gox -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="linux darwin windows" \
	-arch="amd64 386" \
	-osarch="!darwin/arm64" \
	-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .

## Make all distributives
dist: build-all
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf apisite-${VERSION}-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r apisite-${VERSION}-{}.zip {} \; && \
	cd ..

# ------------------------------------------------------------------------------
## Misc

## Count lines of code (including tests) and update LOC.md
cloc: LOC.md

LOC.md: $(SOURCES)
	cloc --by-file --md $(SOURCES) > $@

## List Makefile targets
help:  Makefile
	@grep -A1 "^##" $< | grep -vE '^--$$' | sed -E '/^##/{N;s/^## (.+)\n(.+):(.*)/\t\2:\1/}' | column -t -s ':'
