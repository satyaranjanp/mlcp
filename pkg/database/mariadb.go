package database

import (
	"database/sql"
	"fmt"
	"mlcp/pkg/config"
	_ "github.com/go-sql-driver/mysql"
)

const (
	dbpath string = "mlcp"
)

type MariaDB struct {
	host string
	user string
	password string
}

func initMariaDb() (database, error) {
	return &MariaDB {
		host: config.DatabaseHost,
		user: config.DatabaseUser,
		password: config.DatabasePassword,
	}, nil
}

func (m *MariaDB) Read() (interface{}, error) {
	host := m.user+":"+m.password+"@tcp("+m.host+")/"+dbpath
	db, err := sql.Open("mysql", host)
	if err != nil {
		fmt.Printf("Error connecting to db: %v\n", err)
	}
	defer db.Close()

	return nil, nil
}

func (m *MariaDB) Write(interface{}) error {
/*	host := m.user+":"+m.password+"@tcp("+m.host+")/"+dbpath
	db, err := sql.Open("mysql", host)
	if err != nil {
		fmt.Printf("Error connecting to db: %v\n", err)
	}
	defer db.Close()

	insert, err := db.Query(`INSERT INTO test VALUES ( 2, 'TEST' )`)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

 */
	return nil
}

