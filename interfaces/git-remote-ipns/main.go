package main

import (
	"log"
	"os"

	"github.com/Persper/dgit/core"
)

func main() {

	if len(os.Args) < 3 {
		core.Err_print("Usage: git-remote-ipns remote-name url\n")
		return
	}

	/*
	 * For example: git clone ipns::QmS5mHovjz7soFc7joLu2smafRdNg2QDvBGu4s7EKm29Qv
	 * QmS5mHovjz7soFc7joLu2smafRdNg2QDvBGu4s7EKm29Qv: the ipns key value
	 * IPNS_Key -> IPFS_Key -> Git_Commit_Hash
	 */
	if len(os.Args[2]) > 0 && os.Args[2][0:2] == "Qm" {
		repo_ipfs_hash, _ := core.Transform_ipns_to_ipfs(os.Args[2])
		os.Args[2] = core.Transform_ipfs_to_git(repo_ipfs_hash)
	}

	if err := core.Main(true); err != nil {
		log.Fatal(err)
	}
}
