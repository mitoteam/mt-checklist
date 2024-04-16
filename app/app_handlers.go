package app

import (
	"github.com/mitoteam/goappbase"
)

func DoPreRun() (err error) {
	// open database and migrate schema
	if Db, err = goappbase.DbSchema.Open(); err != nil {
		return err
	}

	//check if root user exists
	if err = checkRootUser(); err != nil {
		return err
	}

	return nil //no errors
}

func DoPostRun() error {
	goappbase.DbSchema.Close()

	return nil //no errors
}
