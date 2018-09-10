SHELL = /bin/sh

.PHONY: default

default: generated.go newfiles.go

generated.go: template codeplugs.json
	go generate

newfiles.go: new.tar.bz2
	go generate

new.tar.bz2: new/*
	cd new && tar cjf ../new.tar.bz2 *
