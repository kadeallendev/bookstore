MAIN_PACKAGE_PATH := ./cmd/bookstore/
BINARY_NAME := bookstore
BINARY_OUTPUT := ./bin/${BINARY_NAME}

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## build: build the application
.PHONY: build
build:
    # Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -o=${BINARY_OUTPUT} ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: build
	/tmp/bin/${BINARY_NAME}
