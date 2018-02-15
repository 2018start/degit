package core

import (
	"github.com/dgraph-io/badger"
	"os"
	"path"
)

//Tracker tracks which hashes are published in IPLD
type Tracker struct {
	kv *badger.DB
}

func NewTracker(gitPath string) (*Tracker, error) {
	ipldDir := path.Join(gitPath, "tracker")
	/*if PathExists(ipldDir) {
		os.RemoveAll(ipldDir)
	}*/
	err := os.MkdirAll(ipldDir, 0755)
	if err != nil {
		return nil, err
	}

	opt := badger.DefaultOptions
	opt.Dir = ipldDir
	opt.ValueDir = ipldDir

	kv, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return &Tracker{
		kv: kv,
	}, nil
}

// Input: refName    Output: hash,err
func (t *Tracker) GetRef(refName string) ([]byte, error) {
	txn := t.kv.NewTransaction(true)
	it, err := txn.Get([]byte(refName))
	if err != nil {
		txn.Commit(nil)
		return nil, err
	}
	str, err2 := it.Value()
	txn.Commit(nil)
	return str, err2
}

// Set RefName and hash
func (t *Tracker) SetRef(refName string, hash []byte) error {
	txn := t.kv.NewTransaction(true)
	err := txn.Set([]byte(refName), hash)
	txn.Commit(nil)
	return err
}

// add entry
func (t *Tracker) AddEntry(hash []byte) error {
	txn := t.kv.NewTransaction(true)
	err := txn.Set(hash, []byte{1})
	txn.Commit(nil)
	return err
}

// check entry
func (t *Tracker) HasEntry(hash []byte) (bool, error) {
	txn := t.kv.NewTransaction(true)
	item, err := txn.Get(hash)
	if err != nil {
		txn.Commit(nil)
		return false, err
	}
	val, err := item.Value()
	txn.Commit(nil)
	//val, err := syncBytes(item.Value)
	return val != nil, err
}

func (t *Tracker) Close() error {
	return t.kv.Close()
}

func syncBytes(get func(func([]byte) error) error) ([]byte, error) {
	var out []byte
	err := get(func(data []byte) error {
		out = data
		return nil
	})

	return out, err
}
