SHELL = /bin/sh

.PHONY: clobber

CODEPLUGDIR = github.com/dalefarnsworth/codeplug
SRCDIR = $(GOPATH)/src/$(CODEPLUGDIR)/genCodeplug
BINDIR = $(GOPATH)/bin
SOURCES = $(SRCDIR)/*.go

$(BINDIR)/genCodeplug: $(SOURCES)
	go install

clobber:
	rm -f $(BINDIR)/genCodeplug
