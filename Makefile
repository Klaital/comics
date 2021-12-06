
gitver:
	GITVER=$(git rev-parse HEAD)
	@echo $(GITVER)

build:
	CGO_ENABLED=0 go build -o comics ./cmd/server
	CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -o ./web/main.wasm ./cmd/weblib

test:
	go test ./...

container: build test gitver
	docker build -t klaital/comics-web:$(GITVER) .

push: container
	docker push klaital/comics-web:$(GITVER)

run: container
	docker-compose -f ./run/docker-compose.local.yml up

run-prod: container
	docker-compose -f ./run/docker-compose.prod.yml up


