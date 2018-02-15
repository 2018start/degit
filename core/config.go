package core

import (
	"log"
	"os"

	"github.com/zpatrick/go-config"
)

/*
 *  The format of the configurate file (the default file: ~/.ipfs/degit.ini)
 *  [ipns]
 *  ipns_alias_name=hash-value
 *
 *  [ipfs]
 *  ipfs_alias_name=hash-value
 */
func Read_IPNS_alias(alias string) (string, error) {
	iniFile_path := os.Getenv("HOME") + "/.ipfs/degit.ini"
	iniFile := config.NewINIFile(iniFile_path)
	c := config.NewConfig([]config.Provider{iniFile})
	if err := c.Load(); err != nil {
		log.Fatal(err)
	}
	find_alias := "ipns." + alias
	hash, find_err := c.String(find_alias)
	if find_err != nil {
		Error_print("Error: can not find " + alias + " from file " + iniFile_path)
		return hash, find_err
	}
	std_print("Transform alias %s to hash %s successfully", alias, hash)
	return hash, find_err
}
