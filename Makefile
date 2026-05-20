BINARY   := wz-top-vip
LDFLAGS  := -trimpath -ldflags="-s -w"

.PHONY: all build build-all clean vet

all: build

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/wz-top-vip

build-all:
	GOOS=linux   GOARCH=amd64  go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64       ./cmd/wz-top-vip
	GOOS=windows GOARCH=amd64  go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe ./cmd/wz-top-vip
	GOOS=darwin  GOARCH=arm64  go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64      ./cmd/wz-top-vip

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
	rm -rf dist/
