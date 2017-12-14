all: deps
	go build ./...

gx:
	go get github.com/whyrusleeping/gx
	go get github.com/whyrusleeping/gx-go

deps: gx
	gx --verbose install --global
	gx-go rewrite

install: 
	go build -o git-remote-ipns ./interfaces/git-remote-ipns/main.go
	mv git-remote-ipns $(GOPATH)/bin
	go build -o git-remote-ipfs ./interfaces/git-remote-ipfs/main.go
	mv git-remote-ipfs $(GOPATH)/bin

clean:
	rm $(GOPATH)/bin/git-remote-ipns $(GOPATH)/bin/git-remote-ipfs

.PHONY: all gx deps install



