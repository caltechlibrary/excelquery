#
# Simple Makefile
#
PROG = xlquery

build: fmt
	go build
	go build ./rss2
	go build -o bin/$(PROG) cmds/$(PROG)/$(PROG).go
	cd webapp && gopherjs build

test: fmt
	go test
	go test ./rss2

fmt: 
	gofmt -w $(PROG).go
	gofmt -w $(PROG)_test.go
	gofmt -w rss2/rss2.go
	gofmt -w rss2/rss2_test.go
	gofmt -w cmds/$(PROG)/$(PROG).go
	gofmt -w webapp/webapp.go
	goimports -w $(PROG).go
	goimports -w $(PROG)_test.go
	goimports -w rss2/rss2.go
	goimports -w rss2/rss2_test.go
	goimports -w cmds/$(PROG)/$(PROG).go
	goimports -w webapp/webapp.go

<<<<<<< HEAD
website:
	./mk-website.bash

save:
=======
save: fmt
	./mk-webapp.bash
>>>>>>> 03706367576b61b5f455b6c1b838d7ee5c5ba8b2
	./mk-website.bash
	git commit -am "quick save"
	git push origin master

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROG)-binary-release.zip ]; then /bin/rm -f $(PROG)-binary-release.zip; fi
	if [ -f webapp/webapp.js ]; then /bin/rm -f webapp/webapp.js; fi

install:
	env GOBIN=$(HOME)/bin go install cmds/$(PROG)/$(PROG).go

release:
	./mk-release.bash

website:
	./mk-webapp.bash
	./mk-website.bash

publish:
	./mk-webapp.bash
	./mk-website.bash
	./publish.bash
