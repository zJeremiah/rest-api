# This specifies the flags to pass to the go binary at runtime for versioning the app.
IMPORT=github.com/rest-api/internal/version
GOFLAGS=-ldflags "-s -w -X ${IMPORT}.Version=${version} -X ${IMPORT}.BuildTimeUTC=`date -u '+%Y-%m-%d_%H:%M:%S'` -X ${IMPORT}.AppName=signal-api"
BLDDIR=deploy/bin

# version is defined by git describe --tags, but can be override if deploying a specific version
ifeq ("${version}", "")
version=$(shell git describe --tags --always)
endif

build: clean
	CGO_ENABLED=0 GOOS=linux go build ${GOFLAGS} -o ${BLDDIR}/
	go build ${GOFLAGS} -o ${BLDDIR}/local/


clean: # clean up build directory
	rm -rf $(BLDDIR)
	mkdir -p $(BLDDIR)
	rm -f ${BLDDIR}/signal-api
	rm -rf ${BLDDIR}/local

test:
	@echo "version: ${version}"
	go test -cover -race ./...

# None of the targets generates a file.
.PHONY: clean test build
