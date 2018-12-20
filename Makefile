ifeq ($(OS),Windows_NT)
  TARGET = pixela_art.exe
else
  TARGET = pixela_art
endif

LDFLAGS := -w -s

SRCS	= \
	src/lib/color/rgb8.go \
	src/lib/date/date.go \
	src/lib/file/file.go \
	src/lib/svg/pixela_svgparser.go \
	src/lib/image/dot_imageparser.go \
	src/setting.go \
	src/main.go

# --------------------------------------------------

all: $(TARGET)

$(TARGET): $(SRCS)
	go get ./...
	go build -ldflags "$(LDFLAGS)" -o $@ ./src

test:
	go test ./src/tests

bench:
	go test ./src/tests -bench .

run:
	./$(TARGET)

clean:
	rm -f $(TARGET)
	
# --------------------------------------------------

.PHONY: clean
