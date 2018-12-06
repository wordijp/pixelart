ifeq ($(OS),Windows_NT)
  TARGET = pixela_art.exe
else
  TARGET = pixela_art
endif

LDFLAGS := -w -s

SRCS	= src/main.go

# --------------------------------------------------

all: $(TARGET)

$(TARGET): $(SRCS)
	go build -ldflags "$(LDFLAGS)" -o $@ ./src

clean:
	rm -f $(TARGET)
	
# --------------------------------------------------

.PHONY: clean
