package test

import (
	"os"
	"path/filepath"
)

// findTestdataDBs test the current directory for the subfolders testdata/dbs. If it doesn't exist go up a directory and
// try again. If it doesn't exist at all, panic.
func FindTestdataDBs() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		if _, err := os.Stat(dir + "/testdata/dbs"); err == nil {
			return dir + "/testdata/dbs"
		}

		if dir == "/" {
			panic("couldn't find testdata/dbs")
		}

		dir, err = filepath.Abs(dir + "/..")
		if err != nil {
			panic(err)
		}
	}
}
