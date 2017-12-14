# git-remote-ipld

[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

> The git-remote-ipld provides interfaces for pushing and pulling commits from/to IPFS!
> This helper is experimental as of now.

TODO: Fill out this long description.

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [Contribute](#contribute)
- [License](#license)

## Background

## Install

### 1. Install Go
The build process for ipfs requires Go 1.8 or higher. If you don't have it: [Download Go 1.8+](https://golang.org/dl/).
Download and decompression: 
```
$ wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz 
  (mac-os: wget https://redirector.gvt1.com/edgedl/go/go1.9.2.darwin-amd64.tar.gz)
$ tar -zxvf go1.8.linux-amd64.tar.gz
$ sudo mv go /usr/local/
```
You'll need to add Go's bin directories to your $PATH environment variable e.g., by adding these lines to your /etc/profile (for a system-wide installation) or $HOME/.profile:
```
$ export GOROOT=/usr/local/go
$ export GOPATH=$HOME/gocode
$ export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
Immediate effect and view go version:
```
$ source /etc/profile
$ go version
```

### 2. Install IPFS
Download:
```
$ go get -u -d github.com/ipfs/go-ipfs
```

Uncomment plugin entries in plugin/loader/preload_list and add plugins into IPFS binary: 
```
$ vim $GOPATH/src/github.com/ipfs/go-ipfs/plugin/loader/preload_list
$ remove '#' in line 6 (#ipldgit github.com/ipfs/go-ipfs/plugin/plugins/git 0)
```

Build and install:
```
$ cd $GOPATH/src/github.com/ipfs/go-ipfs
$ make build
$ make install
```

Test it out:
```
$ ipfs init
initializing IPFS node at /home/hqd/.ipfs
generating 2048-bit RSA keypair...done
......
```

Note: If you have installed ipfs in Linux, you can add the needed plugins that allows augmenting the daemons functionality without recompiling. The detailed can refer to [Plugins.md](https://github.com/ipfs/go-ipfs/blob/master/docs/plugins.md).

### 3. Install dgit
Download dgit and make install (Note: May need to use VPN to download golang repo in CHINA):
```
$ go get github.com/Persper/dgit
$ cd $GOPATH/src/github.com/Persper/dgit
$ make install
```

## API
The dgit support the following interfaces:
```
Clone:
$ git clone ipns::QmULVCL5LGcmKaLMZG1qU6ZZyB8vaL3c5LJtSQsXEu5KKW

Pull:
$ git pull ipns::QmULVCL5LGcmKaLMZG1qU6ZZyB8vaL3c5LJtSQsXEu5KKW

Push:
$ git push ipns::
```

## Usage
#### Just use this repo to test it out:
Use a shell to run ipfs daemon:
```
$ ipfs daemon
```
Clone dgit.git and push it to IPFS; 
```
$ git clone https://github.com/Persper/dgit.git
$ cd dgit
```
Push dgit.git into IPFS:
```
$ git push ipns::
...
Pushed to IPNS as ipns::QmULVCL5LGcmKaLMZG1qU6ZZyB8vaL3c5LJtSQsXEu5KKW

To ipns::
 * [new branch]      master -> master
``` 
Fetch dgit.git from IPFS:
```
$ cd ..
$ git clone ipns::QmULVCL5LGcmKaLMZG1qU6ZZyB8vaL3c5LJtSQsXEu5KKW dgit
```
Now you get the same repo.

Note: Some features like remote tracking are still missing, though the plugin is quite usable.

## Contribute

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT 
