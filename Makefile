TARGET = crescent
COVER = cover.out
SRCS = $(shell find . -name '*.go')

$(TARGET):
	go build -o $(TARGET)

test $(COVER): $(SRCS)
	go test -coverprofile=$(COVER)

cover: $(COVER)
	go tool cover -html=$(COVER)
