build:
	go build

test:
	go test -v ./...

lint:
	go fmt ./...