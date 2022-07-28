# This specifies the flags to pass to the go binary at runtime for versioning the app.
IMPORT=github.com/rest-api/internal/version
GOFLAGS=-ldflags "-s -w -X ${IMPORT}.Version=${version} -X ${IMPORT}.BuildTimeUTC=`date -u '+%Y-%m-%d_%H:%M:%S'` -X ${IMPORT}.AppName=rest-api"
BLDDIR=deploy/bin

# version is defined by git describe --tags, but can be override if deploying a specific version
ifeq ("${version}", "")
version=$(shell git describe --tags --always)
endif

build: clean
	go build ${GOFLAGS} -o ${BLDDIR}

# build the markdown file for slate to build the api documentation
docs: build
	@deploy/bin/rest-api -docs && \
	docker-compose -f "./docs/docker-compose.yml" up;

clean: # clean up build directory
	rm -rf $(BLDDIR)
	mkdir -p $(BLDDIR)

test:
	@echo "version: ${version}"
	go test -cover -race ./...

# None of the targets generates a file.
.PHONY: clean test build
