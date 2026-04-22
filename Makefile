.PHONY: build test clean install

BINARY_NAME=bin/diffq
GO_CMD=go

build:
	$(GO_CMD) build -o $(BINARY_NAME) ./cmd

build-all:
	CGO_ENABLED=0 $(GO_CMD) build -o $(BINARY_NAME) ./cmd

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

run-example-all:\
	run-example \
	run-example1 \
	run-example2 \
	run-example3 \
	run-example4

run-example:
	./$(BINARY_NAME) \
		-d1 "mysql://root:1234567@127.0.0.1:3306/mysql" \
		-q1 "SELECT 1 as id, 'hello' as name" \
		-d2 "mysql://root:1234567@127.0.0.1:3306/mysql" \
		-q2 "SELECT 1 as id, 'hello' as name"


run-example1:
	./$(BINARY_NAME) \
		-d1 "mysql://root:1234567@127.0.0.1:3306/mysql" \
		-q1 "SELECT 1 as id, 'hello' as name union all select 2 as id, 'world' as name" \
		-d2 "mysql://root:1234567@127.0.0.1:3306/mysql" \
		-q2 "SELECT 1 as id, 'hello' as name"


run-example2:
	./$(BINARY_NAME) \
		-d1 "mysql://root:1234567@127.0.0.1:3306/mysql" \
		-d2 "mysql://root:1234567@127.0.0.1:3306/mysql" \
		-q1 "SELECT 1 as id, 'hello' as name" \
		-q2 "SELECT 1 as id, 'hello' as name union all select 2 as id, 'world' as name"

run-example3:
	./$(BINARY_NAME) \
		-d1 "mysql://root:1234567@127.0.0.1:3306/mysql" \
		-d2 "mysql://root:1234567@127.0.0.1:3306/mysql" \
		-q1 "SELECT 1 as id, 'hello' as name union all select 2 as id, 'xxx' as name" \
		-q2 "SELECT 1 as id, 'hello' as name union all select 2 as id, 'world' as name"

run-example4:
	./$(BINARY_NAME) \
		-d1 "mysql://root:1234567@127.0.0.1:3306/mysql" \
		-d2 "mysql://root:1234567@127.0.0.1:3306/mysql" \
		-q1 "SELECT 1 as id, 1.1 as num" \
		-q2 "SELECT 1 as id, 2.2 as num"
