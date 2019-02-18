SHELL = /bin/sh

.PHONY: clobber docfiles

CODEPLUGDIR = github.com/dalefarnsworth/codeplug
SRCDIR = $(GOPATH)/src/$(CODEPLUGDIR)/docCodeplug
JSONFILE = $(GOPATH)/src/$(CODEPLUGDIR)/codeplug/codeplugs.json
BINDIR = $(GOPATH)/bin
SOURCES = $(SRCDIR)/*.go

docfiles: $(BINDIR)/docCodeplug $(JSONFILE)
	$(BINDIR)/docCodeplug $(JSONFILE)

$(BINDIR)/docCodeplug: $(SOURCES)
	go install

clobber:
	rm -f $(BINDIR)/docCodeplug
