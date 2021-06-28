
build:
	CGO_ENABLED=0 go build -o comics ./cmd/server 

test:
	go test ./...

run: test
	go run ./cmd/server

container: build test
	docker build -t klaital/comics-web .

