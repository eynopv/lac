build:
	go build -o lac

test:
	go test ./...

watch:
	./scripts/test_watch.sh

clean:
	go clean
	rm -f lac
