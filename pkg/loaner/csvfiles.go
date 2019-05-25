package loaner

import (
	"log"
	"os"
	"path/filepath"

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

func objsToCSV(fpath string, obj interface{}) error {

	// parse obj to CSV
	out, err := gocsv.MarshalBytes(obj)
	if err != nil {
		return err
	}

	//
	_, err = os.Stat(filepath.Dir(fpath))
	if os.IsNotExist(err) {
		log.Print("is not exist, mkdir ", filepath.Dir(fpath))
		err := os.Mkdir(filepath.Dir(fpath), 0755)
		if err != nil {
			log.Print("mkdir err: ", err)
			return err
		}
	} else if err != nil {
		log.Print("dir stat err: ", err)
		return err
	}
	// open file
	f, err := os.Create(fpath)
	if err != nil {
		log.Print("csv file create err: ", err)
		return err
	}
	defer f.Close()

	_, err = f.Write(out)
	if err != nil {
		log.Print("csv file write err: ", err)
		return err
	}

	return nil
}
