build:
	go build -o bin/tamborete ./cmd/server

run:
	go run ./cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/ data/dump.json