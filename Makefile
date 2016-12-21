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
	goimports -w $(PROJECT).go
	goimports -w $(PROJECT)_test.go
	goimports -w cmds/$(PROJECT)/$(PROJECT).go
	goimports -w webapp/webapp.go

save:
	git commit -am "quick save"
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

release:
	./mk-release.bash

website:
	./mk-webapp.bash
	./mk-website.bash

publish:
	./mk-webapp.bash
	./mk-website.bash
	./publish.bash
