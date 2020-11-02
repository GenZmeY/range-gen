# range-gen creates a list of scene ranges based on a set of frames from the video.
# Copyright (C) 2020 GenZmeY
# mailto: genzmey@gmail.com
#
# This file is part of range-gen.
#
# range-gen is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

NAME     = range-gen
VERSION  = dev_$(shell date +%F_%T)
GOCMD    = go
LDFLAGS := "$(LDFLAGS) -X 'main.Version=$(VERSION)'"
GOBUILD  = $(GOCMD) build -ldflags=$(LDFLAGS)
SRCMAIN  = ./cmd/$(NAME)
BINDIR   = bin
BIN      = $(BINDIR)/$(NAME)
README   = ./doc/README
LICENSE  = LICENSE
PREFIX   = /usr

.PHONY: all prep build doc check-build freebsd-386 darwin-386 linux-386 windows-386 freebsd-amd64 darwin-amd64 linux-amd64 windows-amd64 compile install check-install uninstall clean

all: build

prep: clean
	go mod init; go mod tidy
	mkdir $(BINDIR)
	
build: prep
	$(GOBUILD) -o $(BIN) $(SRCMAIN)
	
doc: check-build
	test -d ./doc || mkdir ./doc
	$(BIN) --help > ./doc/README

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
	
install: check-build doc
	install -m 755 -d                                        $(PREFIX)/bin/
	install -m 755 $(BIN)                                    $(PREFIX)/bin/
	install -m 755 -d                                        $(PREFIX)/share/licenses/$(NAME)/
	install -m 644 $(LICENSE)                                $(PREFIX)/share/licenses/$(NAME)/
	install -m 755 -d                                        $(PREFIX)/share/licenses/$(NAME)/go-perceptualhash
	install -m 644 ./third_party/go-perceptualhash/LICENSE   $(PREFIX)/share/licenses/$(NAME)/go-perceptualhash
	install -m 755 -d                                        $(PREFIX)/share/doc/$(NAME)/
	install -m 644 $(README)                                 $(PREFIX)/share/doc/$(NAME)/

check-install:
	test -e $(PREFIX)/bin/$(NAME) || \
	test -d $(PREFIX)/share/licenses/$(NAME) || \
	test -d $(PREFIX)/share/doc/$(NAME)

uninstall: check-install
	rm -f  $(PREFIX)/bin/$(NAME)
	rm -rf $(PREFIX)/share/licenses/$(NAME)
	rm -rf $(PREFIX)/share/doc/$(NAME)

clean:
	rm -rf $(BINDIR)

