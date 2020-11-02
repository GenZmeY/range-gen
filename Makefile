NAME=range-gen
VERSION=0.1.0
GOCMD=go
LDFLAGS:="$(LDFLAGS) -X 'main.Version=$(VERSION)'"
GOBUILD=$(GOCMD) build -ldflags=$(LDFLAGS)
SRCMAIN=.
BINDIR=bin
BIN=$(BINDIR)/$(NAME)
PREFIX=/usr

all: build

prep: clean
	mkdir $(BINDIR)
	
build: prep
	$(GOBUILD) -o $(BIN) $(SRCMAIN)

check-build:
	test -e $(BIN)

freebsd-386: prep
	GOOS=freebsd GOARCH=386 $(GOBUILD) -o $(BIN)-freebsd-386 $(SRCMAIN)

darwin-386: prep
	GOOS=darwin GOARCH=386 $(GOBUILD) -o $(BIN)-darwin-386 $(SRCMAIN)

linux-386: prep
	GOOS=linux GOARCH=386 $(GOBUILD) -o $(BIN)-linux-386 $(SRCMAIN)

windows-386: prep
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(BIN)-windows-386.exe $(SRCMAIN)
	
freebsd-amd64: prep
	GOOS=freebsd GOARCH=amd64 $(GOBUILD) -o $(BIN)-freebsd-amd64 $(SRCMAIN)

darwin-amd64: prep
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN)-darwin-amd64 $(SRCMAIN)

linux-amd64: prep
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN)-linux-amd64 $(SRCMAIN)

windows-amd64: prep
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN)-windows-amd64.exe $(SRCMAIN)

compile: freebsd-386 darwin-386 linux-386 windows-386 freebsd-amd64 darwin-amd64 linux-amd64 windows-amd64
	
install: check-build
	install -m 755 -d         $(PREFIX)/bin/
	install -m 755 $(BIN)     $(PREFIX)/bin/

check-install:
	test -e $(PREFIX)/bin/$(NAME)

uninstall: check-install
	rm -f  $(PREFIX)/bin/$(NAME)

clean:
	rm -rf $(BINDIR)
