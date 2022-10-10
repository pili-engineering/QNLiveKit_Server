UNAME=$(shell uname)

all: 
	CGO_ENABLED=0 go install -ldflags '-X "github.com/qiniu/version.version=${BUILD_NUMBER} ${BUILD_ID} ${BUILD_URL}" -X github.com/qiniu/version.pkgName=${PACKAGE_NAME}' -v ./app/...

linux: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags '-X "github.com/qiniu/version.version=${BUILD_NUMBER} ${BUILD_ID} ${BUILD_URL}" -X github.com/qiniu/version.pkgName=${PACKAGE_NAME}' -v ./app/...

install: all
	@echo

test:
	go test -race ./... 

coverage:
	go test -race -cover -coverprofile="${QBOXROOT}/coverage.txt" ./...

testv:
	go test -v -race ./...

gofmt-check:
	find . -name "*.go" |  xargs gofmt -s -l -e

govet-check:
	go list ./... | xargs go vet -composites=false
