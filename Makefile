.PHONY: test tests
all: test

test:
	@echo "unit testing..."
	go test -v

tests:
	@echo "run all tests..."
	go test -v . -tags=integration -args ${params}
