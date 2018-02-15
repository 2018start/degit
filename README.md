# degit (decentralized git on IPFS)

[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

> The degit provides interfaces for pushing and pulling commits from/to IPFS/IPNS!
> This helper is experimental as of now.

## Table of Contents

- [Background](#background)
- [Install](#install)
- [API](#api)
- [Usage](#usage)
- [Note](#note)
- [Contribute](#contribute)
- [License](#license)

## Background

## Install

### 1. Install Go
The build process for IPFS requires Go 1.8 or higher. If you don't have it,
[install Go 1.8+](https://golang.org/doc/install).

Here we provide an example script for quick start.

On Ubuntu:
```bash
# Install Go in ~/.local/go
wget https://redirector.gvt1.com/edgedl/go/go1.9.2.linux-amd64.tar.gz
mkdir -p ~/.local
tar -C ~/.local -zxf go1.9.2.linux-amd64.tar.gz

# Add Go's bin directory to the $PATH environment variable
echo 'export GOROOT=$HOME/.local/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOROOT/bin' >> ~/.bashrc
# Set GOPATH which is used later as a workspace for Go
echo 'export GOPATH=$HOME/go' >> ~/.bashrc

# Make the above effective immediately
source ~/.bashrc
go version
```

On macOS:
```bash
# Install Go in /usr/local
curl -LO https://redirector.gvt1.com/edgedl/go/go1.9.2.darwin-amd64.tar.gz
mkdir -p ~/.local
tar -C ~/.local -xf go1.9.2.darwin-amd64.tar.gz

# Add Go's bin directory to the $PATH environment variable
echo 'export GOROOT=$HOME/.local/go' >> ~/.bash_profile
echo 'export PATH=$PATH:$GOROOT/bin' >> ~/.bash_profile
# Set GOPATH which is used later as a workspace for Go
echo 'export GOPATH=$HOME/go' >> ~/.bash_profile

# Make the above effective immediately
source ~/.bash_profile
go version
```

### 2. Build IPFS with IPLD git plugin

Download IPFS and dependencies:
```bash
go get -u -d github.com/ipfs/go-ipfs
```

Uncomment the plugin entry in $GOPATH/src/github.com/ipfs/go-ipfs/plugin/loader/preload_list: `#ipldgit github.com/ipfs/go-ipfs/plugin/plugins/git 0`. Or, directly append the line:
```bash
echo 'ipldgit github.com/ipfs/go-ipfs/plugin/plugins/git 0' >> $GOPATH/src/github.com/ipfs/go-ipfs/plugin/loader/preload_list
```

Build and install:
```bash
cd $GOPATH/src/github.com/ipfs/go-ipfs
make build
make install
```

Test it out:
```bash
ipfs init
```

Note: If you have installed ipfs in Linux, you can add the needed plugin
without recompiling. The detailed can refer to
[Plugins.md](https://github.com/ipfs/go-ipfs/blob/master/docs/plugins.md).

### 3. Install degit

Download degit and make install:
```bash
go get github.com/Persper/degit
make -C $GOPATH/src/github.com/Persper/degit install
```

## API

The degit supports the following interfaces:
```bash
# Clone:
git clone ipns::QmULVCL5LGcmKaLMZG1qU6ZZyB8vaL3c5LJtSQsXEu5KKW 
git clone ipfs::hash-value

# Pull:
git pull ipns::QmULVCL5LGcmKaLMZG1qU6ZZyB8vaL3c5LJtSQsXEu5KKW
git pull ipfs::hash-value

# Push:
git push ipns::
git push ipfs::
```

## Usage
#### Example 1: Create a new repo and store it in IPFS: 
Create a new repo. For example, I create a new repo, named sample.
```bash
$ mkdir sample
$ cd sample
$ git init
```
Then, create a new file (e.g., sample.txt) and commit this added file.
```bash
$ echo "hello world" > sample.txt
$ git add sample.txt
$ git commit -m "hello world"
 [master (root-commit) 74ab5c1] hello world
 1 file changed, 1 insertion(+)
 create mode 100644 sample.txt
```
Next, push this new repo to IPFS.
```bash
$ git push --set-upstream ipns:: master
...
Pushed to IPNS as ipns::QmaqogN63T55e1qxwbFfd2ZVpfsZLhj6ikc1LjYRCi9iP8

To ipns::
 * [new branch]      master -> master
Branch master set up to track remote branch master from ipns::.
```
Set the IPNS::hash as the default remote repo:
```bash
$ git remote add origin ipns::QmaqogN63T55e1qxwbFfd2ZVpfsZLhj6ikc1LjYRCi9iP8
$ git remote -v
origin	ipns::QmaqogN63T55e1qxwbFfd2ZVpfsZLhj6ikc1LjYRCi9iP8 (fetch)
origin	ipns::QmaqogN63T55e1qxwbFfd2ZVpfsZLhj6ikc1LjYRCi9iP8 (push)
$ git branch master -u origin/master
Branch master set up to track remote branch master from origin.
```
Then, you can use "git pull" or "git push" to pull/push new commits from/to the IPNS.

Fetch the sample repo from IPNS:
```bash
$ git clone ipns::QmaqogN63T55e1qxwbFfd2ZVpfsZLhj6ikc1LjYRCi9iP8 sample
```
Now you get the same repo.

#### Example 2: Move an existing repo in Github to IPFS: 
Clone degit.git and push it to IPFS; 
```bash
$ git clone https://github.com/Persper/degit.git
$ cd degit
```
Push degit.git into IPFS:
```bash
$ git push ipns::
...
Pushed to IPNS as ipns::QmTU81e9kr4MeWaLP2gExfqyHSzP3L1wXzymogUNbwxu6P

To ipns::
 * [new branch]      master -> master
```
Set the IPNS::hash as the default remote repo:
```bash
$ git remote set-url origin ipns::QmTU81e9kr4MeWaLP2gExfqyHSzP3L1wXzymogUNbwxu6P
$ git remote -v
origin	ipns::QmTU81e9kr4MeWaLP2gExfqyHSzP3L1wXzymogUNbwxu6P (fetch)
origin	ipns::QmTU81e9kr4MeWaLP2gExfqyHSzP3L1wXzymogUNbwxu6P (push)
```
Then, you can use "git pull" or "git push" to pull/push new commits from/to the IPNS. 

Fetch degit.git from IPNS:
```bash
$ git clone ipns::QmTU81e9kr4MeWaLP2gExfqyHSzP3L1wXzymogUNbwxu6P degit
```
Now you get the same repo.

#### Example 3: Create an alias for the targeted IPNS hash:
Create an alias for the targeted IPNS hash, because the IPNS hash is difficult to remember:
```bash
$ vim ~/.ipfs/degit.ini
[ipns]
sample=QmaqogN63T55e1qxwbFfd2ZVpfsZLhj6ikc1LjYRCi9iP8
```
Then, you can use the alias:
```bash
$ git clone ipns::sample sample
```

## Note

Todo: Some features (e.g., tracking the remote state when issuing git pushes), though the plugin is quite usable.

## Contribute

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT
