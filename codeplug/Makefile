SHELL = /bin/sh

CODEPLUGDIR = github.com/dalefarnsworth/codeplug
SRCDIR = $(GOPATH)/src/$(CODEPLUGDIR)/codeplug

generated.go: $(SRCDIR)/template $(SRCDIR)/codeplugs.json
	go generate
