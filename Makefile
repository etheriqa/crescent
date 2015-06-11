TARGET = crescent
COVER = cover.out
SRCS = $(shell find . -name '*.go')

.PHONY: test cover

$(TARGET): test
	go build -o $(TARGET)

test $(COVER): $(SRCS)
	go test -v -covermode=count -coverprofile=$(COVER)

cover: $(COVER)
	go tool cover -html=$(COVER)
