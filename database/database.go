// Package database takes care of read / write process.
package database

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
)

// Storage describes accesses to a storage
type Storage interface {
	Read() []byte
	Write([]byte) error
}

// File is a storage onto the disk
type File struct{}

// Read fetches data from storage
func (f File) Read() []byte {
	rawJSON, err := os.Open(f.path())
	if err != nil {
		log.Fatal(err)
	}
	defer rawJSON.Close()

	readJSON, err2 := ioutil.ReadAll(rawJSON)
	if err2 != nil {
		log.Fatal(err2)
	}

	return readJSON
}

// Write effectively writes data into storage
func (f File) Write(data []byte) error {
	return ioutil.WriteFile(f.path(), data, 0)
}

func (f File) path() string {
	return os.Getenv("GOPATH") + "/src/github.com/bruno-chavez/restedancestor/database/database.json"
}

// ---------------------------------

// Database describes accesses to a storage
type Database interface {
	Prepare(string, ...interface{}) (Stmt, error)
	LastInsertRowID() int64
}

// Db represents a connection to the storage
type Db struct {
	sqlite *sqlite3.Conn
}

// Prepare encapsulates the inner connection for testability
func (d Db) Prepare(sql string, args ...interface{}) (Stmt, error) {
	return d.sqlite.Prepare(sql, args...)
}

func (d Db) LastInsertRowID() int64 {
	return d.sqlite.LastInsertRowID()
}

// Stmt represents a query statement
// Cf. sqlite3.Stmt
type Stmt interface {
	Close() error
	Step() (bool, error)
	Exec(...interface{}) error
	Scan(dst ...interface{}) error
}

// NewDb initialise a new connection
func NewDb() Database {
	s, err := sqlite3.Open("./database/database.db")
	if err != nil {
		log.Fatal(err)
	}
	s.BusyTimeout(5 * time.Second)

	return Db{
		sqlite: s,
	}
}
