dev:
	encore run

temporal:
# UI at http://localhost:8233
	temporal server start-dev

test:
	go test -v ./...
