SHELL = /bin/sh

.PHONY: default

default: generated.go new.tgz

generated.go: template codeplugs.json
	go generate

new.tgz: new/*
	cd new && tar czf ../new.tgz *
