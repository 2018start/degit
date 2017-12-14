package core

import (
	"log"
	"os"
	"strings"

	cid "github.com/ipfs/go-cid"
)

func Transform_ipfs_to_git(ipfs_hash string) string {
	c2, _ := cid.Decode(ipfs_hash)
	mhash := c2.Hash()
	hash := mhash.HexString()[4:]
	return hash
}

func Transform_ipns_to_ipfs(ipns_key string) (string, error) {
	l := log.New(os.Stderr, "", 0)

	// resolve the ipns key
	l.Printf("\n------ Resolve the key of IPNS: ------")
	s := "ipfs name resolve " + ipns_key
	out, err := exec_shell(s)
	out_arr := strings.Split(out, "/")
	out = out_arr[len(out_arr)-1]
	out = strings.Replace(out, "\n", "", -1)
	l.Printf("The Command of resolving the key: " + s)
	l.Printf("The IPFS hash of the IPNS key: " + out)
	l.Printf("")

	return out, err
}
