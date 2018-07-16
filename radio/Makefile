SHELL = /bin/sh
TARGET = radio
SOURCES = *.go
LIBS = ../dfu/*.go ../userdb/*.go ../stdfu/*.go

.PHONY: default linux windows clean clobber

default: linux windows

linux: $(GOPATH)/bin/$(TARGET)

$(GOPATH)/bin/$(TARGET): $(SOURCES) $(LIBS)
	go install

windows: $(TARGET).exe

$(TARGET).exe: $(SOURCES) $(LIBS)
	GOOS=windows GOARCH=386 go build
