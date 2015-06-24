TARGET = crescent
COVER = cover.out
SRCS = $(shell find . -name '*.go')

.PHONY: all test build cover run watch

all: test build

test:
	go test -v

build: $(TARGET)

cover: $(COVER)
	go tool cover -html=$(COVER)

run: $(TARGET)
	./$(TARGET)

watch:
	fswatch-run --latency 0.1 --include '\.go$$' --exclude '\.' . 'make test'

$(TARGET): $(SRCS)
	go build -o $(TARGET) github.com/etheriqa/crescent/app

$(COVER): $(SRCS)
	go test -v -covermode=count -coverprofile=$(COVER)
