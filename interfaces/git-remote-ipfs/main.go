package main

import (
	"log"
	"os"

	"github.com/Persper/dgit/core"
)

func main() {

	if len(os.Args) < 3 {
		core.Err_print("Usage: git-remote-ipfs remote-name url\n")
		return
	}

	/*
	 * For example: git clone ipfs::QmXAqvdCEnbV6t2VZXu11SykRXZuMYkQtdbUB7tQ16XYES
	 * QmXAqvdCEnbV6t2VZXu11SykRXZuMYkQtdbUB7tQ16XYES: the ipfs key value
	 * IPFS_Key -> Git_Commit_Hash
	 */
	if len(os.Args[2]) > 0 && os.Args[2][0:2] == "Qm" {
		os.Args[2] = core.Transform_ipfs_to_git(os.Args[2])
	}

	if err := core.Main(false); err != nil {
		log.Fatal(err)
	}
}
