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
		log.Print("csv to objs, open file err: ", err)
		return err
	}
	defer f.Close()

	// parse csv
	err = gocsv.UnmarshalFile(f, obj)
	if err != nil {
		log.Print("csv to objs, unmarshal err: ", err)
		return err
	}

	return nil
}

func objsToCSV(fpath string, obj interface{}) error {

	// parse obj to CSV
	out, err := gocsv.MarshalBytes(obj)
	if err != nil {
		log.Print("objs to csv, marshal err: ", err)
		return err
	}

	//
	_, err = os.Stat(filepath.Dir(fpath))
	if os.IsNotExist(err) {
		log.Print("is not exist, mkdir ", filepath.Dir(fpath))
		err := os.Mkdir(filepath.Dir(fpath), 0755)
		if err != nil {
			log.Print("objs to csv, mkdir err: ", err)
			return err
		}
	} else if err != nil {
		log.Print("objs to csv, dir stat err: ", err)
		return err
	}
	// open file
	f, err := os.Create(fpath)
	if err != nil {
		log.Print("objs to csv, create err: ", err)
		return err
	}
	defer f.Close()

	_, err = f.Write(out)
	if err != nil {
		log.Print("objs to csv, write err: ", err)
		return err
	}

	return nil
}
