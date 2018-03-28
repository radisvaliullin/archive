package dbdump

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const (
	archiveFolder = "archive"
)

// Dumper - implements a dumping of db to csv file.
type Dumper struct {
	//
	conf *DumperConf

	// dump goroutine errors chan
	dumpErrChan chan error

	// dump result info
	dumpResInfo chan string
}

// DumperConf -
type DumperConf struct {
	// databases instanses describe.
	Instances []DBConf

	// limit of dump goroutine
	DumpGoLimit int
}

// DBConf - describes db.
type DBConf struct {
	Name     string
	Address  string
	User     string
	Password string
	DBName   string
}

// NewDumper - creates new Dumper.
func NewDumper(conf *DumperConf) *Dumper {
	d := &Dumper{
		conf:        conf,
		dumpErrChan: make(chan error, conf.DumpGoLimit*2),
		dumpResInfo: make(chan string, conf.DumpGoLimit*2),
	}
	return d
}

// ToDump - writes db to files.
func (d *Dumper) ToDump() error {

	// to dump tasks pool
	dumpInstChan := make(chan DBConf, d.conf.DumpGoLimit)

	// limited dump goroutine
	dumpWG := &sync.WaitGroup{}
	for i := 0; i < d.conf.DumpGoLimit; i++ {
		dumpWG.Add(1)
		go func() {
			defer dumpWG.Done()
			for dumpInst := range dumpInstChan {
				d.toDump(dumpInst)
			}
		}()
	}

	// print errors and dump resul info
	printWG := &sync.WaitGroup{}
	printWG.Add(1)
	go func() {
		defer printWG.Done()
		for err := range d.dumpErrChan {
			fmt.Println(err)
		}
	}()
	printWG.Add(1)
	go func() {
		defer printWG.Done()
		for info := range d.dumpResInfo {
			fmt.Println(info)
		}
	}()

	// add tasks to dump
	for _, dumpInst := range d.conf.Instances {
		dumpInstChan <- dumpInst
	}
	close(dumpInstChan)

	// wait of dump goroutine stop
	dumpWG.Wait()

	// wait of print stop
	close(d.dumpErrChan)
	close(d.dumpResInfo)
	printWG.Wait()

	return nil
}

// toDump - dumps db instance to file.
func (d *Dumper) toDump(conf DBConf) {

	//
	tablesInfo := []string{}

	// DB
	dbName := fmt.Sprintf("%v:%v@tcp(%v)/%v", conf.User, conf.Password, conf.Address, conf.DBName)

	// open DB
	db, err := sql.Open("mysql", dbName)
	if err != nil {
		err := fmt.Errorf("mysql db - %v/%v open err %v", conf.Address, conf.DBName, err)
		d.dumpErrChan <- err
		return
	}
	defer db.Close()

	// dump users table
	rnum, err := d.dumpTable(db, conf, "users", []string{"user_id", "name"})
	if err != nil {
		d.dumpErrChan <- err
		return
	}
	tInfo := fmt.Sprintf("db %v - %v/%v from users table dump to users.csv %v lines", conf.Name, conf.Address, conf.DBName, rnum)
	tablesInfo = append(tablesInfo, tInfo)

	// dump sales table
	rnum, err = d.dumpTable(db, conf, "sales", []string{"order_id", "user_id", "order_amount"})
	if err != nil {
		d.dumpErrChan <- err
		return
	}
	tInfo = fmt.Sprintf("db %v - %v/%v from sales table dump to sales.csv %v lines", conf.Name, conf.Address, conf.DBName, rnum)
	tablesInfo = append(tablesInfo, tInfo)

	// Zip
	fpath := filepath.Join(archiveFolder, conf.Name)
	zippath := fpath + ".zip"
	err = FolderToZIP(fpath, zippath, false)
	if err != nil {
		d.dumpErrChan <- err
		return
	}

	//
	d.dumpResInfo <- strings.Join(tablesInfo, "\n")
}

// dumpTable - dumps to file users table.
func (d *Dumper) dumpTable(db *sql.DB, conf DBConf, table string, tbFields []string) (int64, error) {

	var rowsNum int64

	// select query for table fields.
	fields := "*"
	if tbFields != nil {
		fields = strings.Join(tbFields, ", ")
	}
	query := fmt.Sprintf("select %v from %v;", fields, table)
	rows, err := db.Query(query)
	if err != nil {
		return 0, fmt.Errorf("mysql db - %v/%v, select table %v err %v", conf.Address, conf.DBName, table, err)
	}
	defer rows.Close()

	// table rows to file
	dumpDir := filepath.Join(archiveFolder, conf.Name)
	if _, err := os.Stat(dumpDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dumpDir, os.ModePerm); err != nil {
			return 0, fmt.Errorf("mysql db - %v/%v, dump folder create err %v", conf.Address, conf.DBName, err)
		}
	}
	dumpFilePath := filepath.Join(dumpDir, fmt.Sprintf("%v.csv", table))
	dumpFile, err := os.Create(dumpFilePath)
	if err != nil {
		return 0, fmt.Errorf("mysql db - %v/%v, %v.csv open err %v", conf.Address, conf.DBName, table, err)
	}
	defer dumpFile.Close()
	dumpCSV := csv.NewWriter(dumpFile)

	// CSV Header
	cols, err := rows.Columns()
	if err != nil {
		return 0, fmt.Errorf("mysql db - %v/%v, %v table columns read err %v", conf.Address, conf.DBName, table, err)
	}
	if err := dumpCSV.Write(cols); err != nil {
		return 0, fmt.Errorf("mysql db - %v/%v, %v table columns write to csv err %v", conf.Address, conf.DBName, table, err)
	}
	dumpCSV.Flush()

	// scan rows and write to csv.
	rawVals := make([][]byte, len(cols))
	vals := make([]string, len(cols))
	dest := make([]interface{}, len(cols))
	for i := range rawVals {
		dest[i] = &rawVals[i]
	}
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return 0, fmt.Errorf("mysql db - %v/%v, %v table rows scan err %v", conf.Address, conf.DBName, table, err)
		}
		rowsNum++

		for i, rVal := range rawVals {
			if rVal == nil {
				vals[i] = "null"
			} else {
				vals[i] = string(rVal)
			}
		}

		if err := dumpCSV.Write(vals); err != nil {
			return 0, fmt.Errorf("mysql db - %v/%v, %v table values write err %v", conf.Address, conf.DBName, table, err)
		}
		dumpCSV.Flush()
	}
	return rowsNum, nil
}
