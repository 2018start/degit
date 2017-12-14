all: deps
	go build ./...

gx:
	go get github.com/whyrusleeping/gx
	go get github.com/whyrusleeping/gx-go

deps: gx
	gx --verbose install --global
	gx-go rewrite

install: 
	go build -o git-remote-ipns 
	mv git-remote-ipns $(GOPATH)/bin

.PHONY: all gx deps install
