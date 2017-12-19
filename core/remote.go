package core

import (
	"os"
	"path"
)

func fetch_remote_repo(gitPath string) (string, error) {
	remoteDir := path.Join(gitPath, "remote-repo")
	if PathExists(remoteDir) {
		os.RemoveAll(remoteDir)
	}

	err := os.MkdirAll(remoteDir, 0755)
	if err != nil {
		return remoteDir, err
	}

	return remoteDir, nil
}
