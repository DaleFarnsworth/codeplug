SHELL = /bin/sh

.PHONY: default

default: generated.go newfiles.go

generated.go: template codeplugs.json
	go generate

newfiles.go: new.tgz
	go generate

new.tgz: new/*
	cd new && tar czf ../new.tgz *
