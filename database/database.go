package database

import (
	"io/ioutil"
	"log"
	"os"
)

type Database struct{}

// read fetches data from storage
func (d Database) read() []byte {
	rawJSON, err := os.Open(path())
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

// write effectively writes data into storage
func (d Database) write(data []byte) error {
	return ioutil.WriteFile(path(), data, 0)
}

func path() string {
	return os.Getenv("GOPATH") + "/src/github.com/bruno-chavez/restedancestor/database/database.json"
}
