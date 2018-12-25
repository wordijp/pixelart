ifeq ($(OS),Windows_NT)
  TARGET = example.exe
else
  TARGET = example
endif

LDFLAGS := -w -s

SRCS	= main.go

# --------------------------------------------------

all: $(TARGET)

$(TARGET): $(SRCS)
	go get ./...
	go build -ldflags "$(LDFLAGS)" -o $@ .

run:
	./$(TARGET)

clean:
	rm -f $(TARGET)
	
# --------------------------------------------------

.PHONY: clean
