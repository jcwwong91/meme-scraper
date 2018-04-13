
GO_URL=https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz
GO_TAR=$(shell basename $(GO_URL))
pwd=$(shell pwd)
GOPATH=$(pwd)

all:: install

tools::
	wget $(GO_URL)
	tar -xzvf $(GO_TAR)
	mkdir -p src/pkg/bin
	$(GOPATH)/go/bin/go get golang.org/x/net/html
	$(GOPATH)/go/bin/go get github.com/anaskhan96/soup
	

install::
	$(GOPATH)/go/bin/go install meme-scraper 


run::
	$(GOPATH)/bin/meme-scraper
