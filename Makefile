.PHONY: test tests
all: test

test:
	@echo "unit testing..."
	go test -v -timeout 30s

tests:
	@echo "run all tests..."
	go test -v -timeout 30s . -tags=integration -args ${params}
