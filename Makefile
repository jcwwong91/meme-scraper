
GO_URL=https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz
GO_TAR=$(shell basename $(GO_URL))

tools::
	wget $(GO_URL)
	tar -xzvf $(GO_TAR)
	mkdir -p src/pkg/bin
	

install:
	export GOPATH
