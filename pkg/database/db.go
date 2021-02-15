package database

import (
	"github.com/golang/glog"
	"mlcp/pkg/common"
	"mlcp/pkg/config"
	"time"
)

type DbType string

const (
	MariaDb DbType = "mariadb"
	// write to database buffer count
	writeCount int = 5
	// wait time before next writing of writeCount number of data
	writeCoolDownPeriod = 5*time.Second
)

type SlotData struct {
	SlotId uint32
	SlotType common.SlotType
}

type VehicleData struct {
	RegnNo string
	Type string
	InTime time.Time
	OutTime time.Time
}

// write to database buffer
// this is to handle uncertainty of data writes into database in case database failures
var dataToWrite []interface{}

func InitializeDatabase() (*Database, error) {
	dataToWrite = make([]interface{}, writeCount)
	db, err := initDb()
	if err != nil {
		return nil, err
	}
	go func(db *Database) {
		for ;; {
			if l := len(dataToWrite); l >= writeCount {
				// send writeCount numbers of data to db in a loop
				// if more than writeCount data is present in dataToWrite slice, then those will be handled in the next loop
				cnt := 0
				for ;; {
					if err := db.database.Write(dataToWrite[0]); err != nil {
						glog.Errorf("Error inserting data (%v) to db: %v; will try again.", dataToWrite[0], err)
						break
					}
					// Remove the successful written data from the beginning of slice to preserve memory
					dataToWrite = dataToWrite[1:]
					cnt++
					if cnt == writeCount {
						break
					}
				}
			}
			time.Sleep(writeCoolDownPeriod)
		}
	}(db)
	return db, nil
}

type database interface{
	Read() (interface{}, error)
	Write(data interface{}) error
}

type Database struct {
	database
}

func initDb() (*Database, error) {
	var db database
	var err error
	switch config.DatabaseDriver {
	case string(MariaDb):
		 db, err = initMariaDb()
	}

	return &Database{db}, err
}

func (db *Database) Read() (interface{}, error) {
	return db.database.Read()
}

func (db *Database) Write(data ...interface{}) {
	dataToWrite = append(dataToWrite, data...)
}