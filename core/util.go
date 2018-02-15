package core

import (
	"bytes"
	"compress/zlib"
	"log"
	"os"
	"path"
)

func Error_print(out string) {
	l := log.New(os.Stderr, "", 0)
	l.Printf(out)
}

func std_print(format string, a ...interface{}) {
	l := log.New(os.Stderr, "", 0)
	l.Printf(format, a...)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func getLocalDir() (string, error) {
	localdir := path.Join(os.Getenv("GIT_DIR"))

	if err := os.MkdirAll(localdir, 0755); err != nil {
		return "", err
	}

	return localdir, nil
}

func compressObject(in []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(in)
	w.Close()
	return b.Bytes()
}
