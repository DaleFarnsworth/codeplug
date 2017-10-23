SHELL = /bin/sh

.PHONY: default

default: generated.go new.tgz

generated.go: $(SRCDIR)/template $(SRCDIR)/codeplugs.json
	go generate

new.tgz: new/*
	cd new && tar czf ../new.tgz *
