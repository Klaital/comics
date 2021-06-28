
build:
	CGO_ENABLED=0 go build -o comics ./cmd/server
	CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -o ./web/main.wasm ./cmd/weblib

test:
	go test ./...

run: test
	go run ./cmd/server

container: build test
	docker build -t klaital/comics-web .

