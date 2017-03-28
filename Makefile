#
# Simple Makefile
#
PROJECT = excelquery

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\"  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

build:
	env CGO_ENABLED=0 go build -o bin/$(PROJECT) cmds/$(PROJECT)/$(PROJECT).go
	cd webapp && gopherjs build

test:
	go test

fmt: 
	gofmt -w $(PROJECT).go
	gofmt -w $(PROJECT)_test.go
	gofmt -w cmds/$(PROJECT)/$(PROJECT).go
	gofmt -w webapp/webapp.go

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "quick save"; fi
	git push origin $(BRANCH)

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f webapp/webapp.js ]; then /bin/rm -f webapp/webapp.js; fi
	if [ -f webapp/webapp.js.map ]; then /bin/rm -f webapp/webapp.js.map; fi
	if [ -f webapp/index.html ]; then /bin/rm -f webapp/index.html; fi
	if [ -f "$(PROJECT)-$(VERSION)-release.zip" ]; then /bin/rm -f "$(PROJECT)-$(VERSION)-release.zip"; fi

install:
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/$(PROJECT)/$(PROJECT).go

webapp:
	./mk-webapp.bash

website:
	./mk-webapp.bash
	./mk-website.bash

publish:
	./mk-webapp.bash
	./mk-website.bash
	./publish.bash

release: dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7
	mkdir -p dist
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v excelquery.md dist/
	zip -r $(PROJECT)-$(VERSION)-release.zip dist/*

dist/linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/excelquery cmds/excelquery/excelquery.go

dist/windows-amd64:
	env GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/excelquery.exe cmds/excelquery/excelquery.go

dist/macosx-amd64:
	env GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/excelquery cmds/excelquery/excelquery.go

dist/raspbian-arm7:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/excelquery cmds/excelquery/excelquery.go



