#!/bin/bash

run_tests() {
	go test ./...
}

run_tests

while true; do
	inotifywait -r -e modify,move,create,delete --include '.*\.go$' . |
		while read -r file; do
				run_tests
		done
done
