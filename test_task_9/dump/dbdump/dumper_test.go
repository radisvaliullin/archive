package dbdump

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// TestDumper -
func TestDumper(t *testing.T) {

	var (
		country0 = "testcountry0"
		country1 = "testcountry1"
		country2 = "testcountry2"
	)

	dbaddress := os.Getenv("mysql_addr")
	if dbaddress == "" {
		t.Fatal("set env mysql_addr")
	}
	dbuser := os.Getenv("mysql_user")
	if dbuser == "" {
		t.Fatal("set env mysql_user")
	}
	dbpwd := os.Getenv("mysql_pwd")
	if dbpwd == "" {
		t.Fatal("set env mysql_pwd")
	}

	// open DB
	dbName := fmt.Sprintf("%v:%v@tcp(%v)/", dbuser, dbpwd, dbaddress)
	db, err := sql.Open("mysql", dbName)
	if err != nil {
		t.Fatalf("test db open err %v", err)
	}
	defer db.Close()

	// create test databases
	testDatabases := []string{country0, country1, country2}

	for _, tdb := range testDatabases {
		_, err = db.Exec("CREATE DATABASE " + tdb)
		if err != nil {
			t.Fatalf("test db %v create err %v", tdb, err)
		}

		_, err = db.Exec("USE " + tdb)
		if err != nil {
			t.Fatalf("test db %v use err %v", tdb, err)
		}

		createTable := "CREATE TABLE users (" +
			" user_id INT unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			" name varchar(255)" +
			" );"
		_, err = db.Exec(createTable)
		if err != nil {
			t.Fatalf("test db %v users table create err %v", tdb, err)
		}

		insertValues := "INSERT INTO users (name) VALUES ('Jose'), ('Juan'), ('Miguel');"
		_, err = db.Exec(insertValues)
		if err != nil {
			t.Fatalf("test db %v users table insert values err %v", tdb, err)
		}

		createTable = "CREATE TABLE sales (" +
			" order_id INT unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			" user_id INT unsigned NOT NULL," +
			" order_amount FLOAT NOT NULL" +
			" );"
		_, err = db.Exec(createTable)
		if err != nil {
			t.Fatalf("test db %v sales table create err %v", tdb, err)
		}

		insertValues = "INSERT INTO sales (user_id, order_amount) VALUES (1, 11.0), (2, 22.0), (3, 37.73);"
		_, err = db.Exec(insertValues)
		if err != nil {
			t.Fatalf("test db %v users table insert values err %v", tdb, err)
		}

	}

	// test databases to dump
	dconf := &DumperConf{
		Instances: []DBConf{
			{
				Name:     country0,
				Address:  dbaddress,
				User:     dbuser,
				Password: dbpwd,
				DBName:   country0,
			},
			{
				Name:     country1,
				Address:  dbaddress,
				User:     dbuser,
				Password: dbpwd,
				DBName:   country1,
			},
			{
				Name:     country2,
				Address:  dbaddress,
				User:     dbuser,
				Password: dbpwd,
				DBName:   country2,
			},
		},
		DumpGoLimit: 2,
	}
	d := NewDumper(dconf)
	d.ToDump()

	// drop test databases
	for _, tdb := range testDatabases {
		_, err = db.Exec("DROP DATABASE " + tdb)
		if err != nil {
			t.Fatalf("test db %v drop err %v", tdb, err)
		}
	}

	t.Log("TEST END")
}
