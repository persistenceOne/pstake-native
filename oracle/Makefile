hello:
	echo


GOBIN = $(shell go env GOPATH)/bin
GOARCH = $(shell go env GOARCH)
GOOS = $(shell go env GOOS)


verify:
	@echo "verifying modules"
	@go mod verify

all:


install:
ifeq (${OS},Windows_NT)
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/orc.exe ./

else
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/orc ./node

endif

build:
ifeq (${OS},Windows_NT)
	go build  ${BUILD_FLAGS} -o build/${GOOS}/${GOARCH}/persistenceCore.exe ./node

else
	go build  ${BUILD_FLAGS} -o build/${GOOS}/${GOARCH}/persistenceCore ./node

endif
