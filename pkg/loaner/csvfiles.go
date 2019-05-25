package loaner

import (
	"os"

	"github.com/gocarina/gocsv"
)

func csvToObjs(fpath string, obj interface{}) error {

	// open file
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	// parse csv
	err = gocsv.UnmarshalFile(f, obj)
	if err != nil {
		return err
	}

	return nil
}
