.PHONY: test units.test
all: test

units.test:
	go test -timeout 30s .

test: units.test
	go test -timeout 30s . -tags=integration -args ${PARAMS}
