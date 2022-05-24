.PHONY: test
all: test

test:
	@echo "unit testing..."
	go test -v ./...
