package core

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"strings"

	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func Main(use_ipns bool) error {

	printf := func(format string, a ...interface{}) (n int, err error) {
		return fmt.Printf(format, a...)
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

		switch {
		case command == "capabilities":
			printf("push\n")
			printf("fetch\n")
			printf("\n")
		case strings.HasPrefix(command, "list"):

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
					return err
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

			if use_ipns == true {
				exec_publish(c.String())
			} else {
				std_print("Pushed to IPFS as \x1b[32mipld::%s\x1b[39m\n", headHash)
				std_print("Head CID is %s\n", c.String())
			}

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