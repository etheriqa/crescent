TARGET = crescent
COVER = cover.out
SRCS = $(shell find . -name '*.go')

.PHONY: test cover watch

$(TARGET): test
	go build -o $(TARGET)

test:
	go test -v

$(COVER): $(SRCS)
	go test -v -covermode=count -coverprofile=$(COVER)

cover: $(COVER)
	go tool cover -html=$(COVER)

watch:
	fswatch -o $(SRCS) | while read line; do clear; date; echo; make; done
