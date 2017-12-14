package core

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	//"strconv"
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
