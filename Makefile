TARGET = crescent
COVER = cover.out
SRCS = $(shell find . -name '*.go')
TARGET_PACKAGE = github.com/etheriqa/crescent
TEST_PACKAGE = github.com/etheriqa/crescent/game

.PHONY: all test build cover clean watch

all: test build

test:
	go test -v $(TEST_PACKAGE)

build: $(TARGET)

cover: $(COVER)
	go tool cover -html=$(COVER)

clean:
	rm -f $(TARGET) $(COVER)

watch:
	fswatch-run --latency 0.1 --include '\.go$$' --exclude '\.' . 'make test'

$(TARGET): $(SRCS)
	go build -o $(TARGET) $(TARGET_PACKAGE)

$(COVER): $(SRCS)
	go test -v -covermode=count -coverprofile=$(COVER) $(TEST_PACKAGE)
