SHELL = /bin/sh
TARGET = radio
SOURCES = *.go
LIBS = ../dfu/*.go ../userdb/*.go ../stdfu/*.go

.PHONY: default linux windows clean clobber

default: windows

linux: $(GOPATH)/bin/$(TARGET)

$(GOPATH)/bin/$(TARGET): $(SOURCES) $(LIBS)
	go install

windows: /home/dale/$(TARGET).exe

/home/dale/$(TARGET).exe: $(TARGET).exe
	cp $(TARGET).exe /home/dale/windows/

$(TARGET).exe: $(SOURCES) $(LIBS)
	GOOS=windows GOARCH=386 go build
