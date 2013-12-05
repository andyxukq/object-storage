# Makefile that builds and installs all sub packages (located under ./...)
# Currently only useful to get automatic version information included in the build.

VERSIONCMD = "`git symbolic-ref HEAD | cut -b 12-`-`git rev-parse HEAD`"
VERSION = $(shell echo $(VERSIONCMD))
DATE = $(shell echo `date +%FT%T%z`)

.PHONY: all build clean test

all: clean build

build:
	export GOPATH="$(CURDIR)/../.." && \
	go install  -tags="no_development" -ldflags "-X tools/version.version $(VERSION) -X tools/version.date $(DATE)" ./...

clean:
	-export GOPATH="$(CURDIR)/../.." && \
	rm -r $$GOPATH/pkg/* || true        
