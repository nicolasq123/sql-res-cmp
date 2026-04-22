.PHONY: build test clean install

BINARY_NAME=diffq
GO_CMD=go

build:
	$(GO_CMD) build -o $(BINARY_NAME) ./cmd

test:
	$(GO_CMD) test ./... -v -cover

test-coverage:
	$(GO_CMD) test ./... -coverprofile=coverage.out
	$(GO_CMD) tool cover -html=coverage.out -o coverage.html

clean:
	rm -f $(BINARY_NAME) coverage.out coverage.html

install: build
	cp $(BINARY_NAME) /usr/local/bin/

lint:
	$(GO_CMD) vet ./...

fmt:
	$(GO_CMD) fmt ./...

run-example:
	./$(BINARY_NAME) \
		-d1 "mysql://root:123456@127.0.0.1:3306/mysql" \
		-q1 "SELECT 1 as id, 'hello' as name" \
		-d2 "mysql://root:123456@127.0.0.1:3306/mysql" \
		-q2 "SELECT 1 as id, 'hello' as name"
