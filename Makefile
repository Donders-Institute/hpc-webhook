ifndef GOPATH
	GOPATH := $(HOME)/go
endif

ifndef GOOS
	GOOS := linux
endif

ifndef GO111MODULE
	GO111MODULE := on
endif

all: build

build:
	GOPATH=$(GOPATH) GOOS=$(GOOS) go build -o bin/hpc-webhook-server cmd/server/server.go

test:
	@GOPATH=$(GOPATH) GOOS=$(GOOS) GOCACHE=off go test -v github.com/Donders-Institute/hpc-webhook/...

clean:
	@rm -rf bin
