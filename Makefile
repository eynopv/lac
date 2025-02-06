build:
	go build -o lac

test:
	go test ./...

clean:
	go clean
	rm -f lac
