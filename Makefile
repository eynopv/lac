build:
	go build -o lac

test:
	go test ./...

lint:
	golangci-lint run

clean:
	go clean
	rm -f lac
