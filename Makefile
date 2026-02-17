.PHONY: build build-all build-linux build-windows build-macos clean

BINARY_NAME=ipbak
CMD_PATH=./cmd/ipbak

build:
	go build -o $(BINARY_NAME) $(CMD_PATH)

build-all: build-linux build-windows build-macos

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-macos-amd64 $(CMD_PATH)
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-macos-arm64 $(CMD_PATH)

clean:
	rm -rf bin/
	rm -f $(BINARY_NAME)
