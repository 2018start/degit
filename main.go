package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	//"strconv"
	"strings"

	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func exec_shell(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		l := log.New(os.Stderr, "", 0)
		l.Printf("exec_shell error: " + s)
	}

	return out.String(), err
}

func exec_shell_may_error(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		//l := log.New(os.Stderr, "", 0)
		//l.Printf("exec_shell_may_error error: " + s)
	}

	return out.String(), err
}

func exec_publish(git_hash string) (string, error) {
	l := log.New(os.Stderr, "", 0)

	// get the name of the targeted key from the repo name
	dir, err := exec_shell("pwd")
	dir_arr := strings.Split(dir, "/")
	key := dir_arr[len(dir_arr)-1]

	// generate the key of ipns
	l.Printf("\n------ Generate the key of IPNS: ------")
	s := "ipfs key gen --type=rsa --size=1024 " + key
	out, err := exec_shell_may_error(s)
	l.Printf("Generate the key: " + s)
	l.Printf("The IPNS key: " + out)

	// publish the hash to the targeted key using IPNS
	l.Printf("\n------ Publish the hash to the key: ------")
	key = strings.Replace(key, "\n", "", -1)
	s = "ipfs name publish --key=" + key + " " + git_hash + " --ttl=8760h"
	out, err = exec_shell(s)
	l.Printf("Publish the hash of the git commit: " + s)
	l.Printf("Pushed to IPNS as \x1b[32mipns::%s\x1b[39m\n\n", out[13:59])

	return out[13:59], err
}

func exec_resolve(ipns_key string) (string, error) {
	l := log.New(os.Stderr, "", 0)

	// resolve the ipns key
	l.Printf("\n------ Resolve the key of IPNS: ------")
	s := "ipfs name resolve " + ipns_key
	out, err := exec_shell(s)
	out_arr := strings.Split(out, "/")
	out = out_arr[len(out_arr)-1]
	out = strings.Replace(out, "\n", "", -1)
	l.Printf("Resolve the key: " + s)
	l.Printf("The IPFS hash of the IPNS key: " + out)
	l.Printf("")

	return out, err
}

func transform_ipfs_to_git(ipfs_hash string) string {
	c2, _ := cid.Decode(ipfs_hash)
	mhash := c2.Hash()
	hash := mhash.HexString()[4:]
	return hash
}

func getLocalDir() (string, error) {
	localdir := path.Join(os.Getenv("GIT_DIR"))

	if err := os.MkdirAll(localdir, 0755); err != nil {
		return "", err
	}

	return localdir, nil
}

func Main() error {
	//l := log.New(os.Stderr, "", 0)
	//l.Println(os.Args)

	printf := func(format string, a ...interface{}) (n int, err error) {
		return fmt.Printf(format, a...)
	}

	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: git-remote-ipld remote-name url")
	}

	localDir, err := getLocalDir()
	if err != nil {
		return err
	}

	repo, err := git.PlainOpen(localDir)
	if err == git.ErrWorktreeNotProvided {
		repoRoot, _ := path.Split(localDir)

		repo, err = git.PlainOpen(repoRoot)
		if err != nil {
			return err
		}
	}

	tracker, err := NewTracker(localDir)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}
	defer tracker.Close()

	stdinReader := bufio.NewReader(os.Stdin)

	for {
		command, err := stdinReader.ReadString('\n')
		if err != nil {
			return err
		}

		command = strings.Trim(command, "\n")

		//l.Printf("< %s", command)
		switch {
		case command == "capabilities":
			printf("push\n")
			printf("fetch\n")
			printf("\n")
		case strings.HasPrefix(command, "list"):
			// git clone ipns::QmS5mHovjz7soFc7joLu2smafRdNg2QDvBGu4s7EKm29Qv
			// QmS5mHovjz7soFc7joLu2smafRdNg2QDvBGu4s7EKm29Qv: the ipns key value
			// IPNS_Key -> IPFS_Key -> Git_Commit_Hash
			if len(os.Args[2]) > 0 && os.Args[2][0:2] == "Qm" {
				repo_ipfs_hash, _ := exec_resolve(os.Args[2])
				os.Args[2] = transform_ipfs_to_git(repo_ipfs_hash)
			}

			headRef, err := repo.Reference(plumbing.HEAD, false)
			if err != nil {
				return err
			}

			it, err := repo.Branches()
			if err != nil {
				return err
			}

			var n int
			err = it.ForEach(func(ref *plumbing.Reference) error {
				n++
				r, err := tracker.GetRef(ref.Name().String())
				if err != nil {
					//return err
				}
				if r == nil {
					r = make([]byte, 20)
				}

				if !strings.HasPrefix(command, "list for-push") && headRef.Target() == ref.Name() && headRef.Type() == plumbing.SymbolicReference && len(os.Args) >= 3 {
					sha, err := hex.DecodeString(os.Args[2])
					if err != nil {
						return err
					}
					if len(sha) != 20 {
						return errors.New("invalid hash length")
					}

					printf("%s %s\n", os.Args[2], headRef.Target().String())
				} else {
					printf("%s %s\n", hex.EncodeToString(r), ref.Name())
				}

				return nil
			})
			it.Close()
			if err != nil {
				return err
			}

			if n == 0 && !strings.HasPrefix(command, "list for-push") && len(os.Args) >= 3 {
				sha, err := hex.DecodeString(os.Args[2])
				if err != nil {
					return err
				}
				if len(sha) != 20 {
					return errors.New("invalid hash length")
				}

				printf("%s %s\n", os.Args[2], "refs/heads/master")
			}

			switch headRef.Type() {
			case plumbing.HashReference:
				printf("%s %s\n", headRef.Hash(), headRef.Name())
			case plumbing.SymbolicReference:
				printf("@%s %s\n", headRef.Target().String(), headRef.Name())
			}

			printf("\n")
		case strings.HasPrefix(command, "push "):
			refs := strings.Split(command[5:], ":")

			localRef, err := repo.Reference(plumbing.ReferenceName(refs[0]), true)
			if err != nil {
				return fmt.Errorf("command push: %v", err)
			}

			headHash := localRef.Hash().String()

			push := NewPush(localDir, tracker, repo)
			err = push.PushHash(headHash)
			if err != nil {
				return fmt.Errorf("command push: %v", err)
			}

			hash := localRef.Hash()
			tracker.SetRef(refs[1], (&hash)[:])

			mhash, err := mh.FromHexString("1114" + headHash)
			if err != nil {
				return fmt.Errorf("fetch: %v", err)
			}

			c := cid.NewCidV1(cid.GitRaw, mhash)

			// IPFS_Hash -> IPNS_Key
			exec_publish(c.String())
			/*l.Printf("Pushed to IPFS as \x1b[32mipld::%s\x1b[39m\n", headHash)
			l.Printf("Head CID is %s\n", c.String())*/
			printf("ok %s\n", refs[0])
			printf("\n")
		case strings.HasPrefix(command, "fetch "):
			parts := strings.Split(command, " ")

			fetch := NewFetch(localDir, tracker)
			err := fetch.FetchHash(parts[1])
			if err != nil {
				return fmt.Errorf("command fetch: %v", err)
			}

			sha, err := hex.DecodeString(parts[1])
			if err != nil {
				return fmt.Errorf("push: %v", err)
			}

			tracker.SetRef(parts[2], sha)

			printf("\n")
		case command == "\n":
			return nil
		case command == "":
			return nil
		default:
			return fmt.Errorf("Received unknown command %q", command)
		}
	}
	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
