.PHONY: release

none:
	@echo "Please use make (build|deps|release)"

deps:
	dep ensure

build:
	@go generate main.go
	@mkdir -p release/
	@GOARCH=amd64 GOOS=linux 	go build -o release/howe-linux-amd64 -ldflags="-s -w" 	main.go
	@GOARCH=amd64 GOOS=darwin 	go build -o release/howe-darwin-amd64 -ldflags="-s -w" 	main.go

release:
	@upx --brute release/howe-linux-amd64
	@upx --brute release/howe-darwin-amd64
	@shasum release/* > release/shasums
