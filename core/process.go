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

func do_put_check(localDir string, hash string, ref_name string) error {

	_, err := exec_shell_may_error("git cat-file -p " + hash)
	//std_print("out: %s\n", out)
	if err != nil {
		std_print("error: failed to push some refs to remote IPNFS repo. \n")
		std_print("hint: Updates were rejected because the remote contains work that you do not have locally. \n")
		std_print("This is usually caused by another repository pushing to the same ref. \n")
		std_print("You may want to first integrate the remote changes (e.g., 'git pull ...') before pushing again. \n\n")
		return err
	}

	/*localDir = path.Join(localDir, "objects")
	localDir = path.Join(localDir, hash[:2])
	localDir = path.Join(localDir, hash[2:])
	std_print("push-check hash: %s ;  path: %s\n", hash, localDir)

	_, err := os.Stat(localDir)
	if os.IsNotExist(err) {
		std_print("error: failed to push some refs to remote IPNFS repo. \n")
		std_print("hint: Updates were rejected because the remote contains work that you do not have locally. \n")
		std_print("This is usually caused by another repository pushing to the same ref. \n")
		std_print("You may want to first integrate the remote changes (e.g., 'git pull ...') before pushing again. \n\n")
		return err
	}*/

	return nil
}

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

		//std_print("command:%s", command)
		//std_print("lensize:%d", len(command))
		switch {
		case command == "capabilities":
			printf("push\n")
			printf("fetch\n")
			printf("\n")
		case strings.HasPrefix(command, "list"): // list for-push ;  list
			//remote_dir, _ := fetch_remote_repo(localDir)
			//std_print(remote_dir)

			/* generate ref: refs/heads/master HEAD. */
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

					//std_print("1 %s %s\n", os.Args[2], headRef.Target().String())
					printf("%s %s\n", os.Args[2], headRef.Target().String())
				} else {
					//std_print("2 %s %s\n", hex.EncodeToString(r), ref.Name())
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

				//std_print("3 %s %s\n", os.Args[2], "refs/heads/master")
				printf("%s %s\n", os.Args[2], "refs/heads/master")
			}

			switch headRef.Type() {
			case plumbing.HashReference:
				//std_print("4 %s %s\n", headRef.Hash(), headRef.Name())
				printf("%s %s\n", headRef.Hash(), headRef.Name())
			case plumbing.SymbolicReference:
				//std_print("5 %s %s\n", headRef.Target().String(), headRef.Name())
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

			std_print("check hash=%s  ; refs[0]=%s\n", os.Args[2], refs[0])
			if os.Args[2] != "" {
				err = do_put_check(localDir, os.Args[2], refs[0])
				if err != nil {
					return err
				}
			}

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
				std_print("Pushed to IPFS as \x1b[32mipfs::%s\x1b[39m\n", headHash)
				std_print("Head CID is %s\n", c.String())
			}

			//std_print("refs[0]=%s\n", refs[0])
			printf("ok %s\n", refs[0])
			printf("\n")
		case strings.HasPrefix(command, "fetch "):
			parts := strings.Split(command, " ")

			//std_print("parts[1]=%s parts[2]=%s\n", parts[1], parts[2])

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
