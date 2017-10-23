package db

import (
	"errors"
	"flag"
	"fmt"

	m_user "github.com/laidingqing/dabanshan/svcs/user/model"
)

// Database represents a simple interface so we can switch to a new system easily
type Database interface {
	Init() error
	GetUserByName(string) (m_user.User, error)
	GetUser(string) (m_user.User, error)
}

var (
	database string
	//DefaultDb is the database set for the microservice
	DefaultDb Database
	//DBTypes is a map of DB interfaces that can be used for this service
	DBTypes = map[string]Database{}
	//ErrNoDatabaseFound error returnes when database interface does not exists in DBTypes
	ErrNoDatabaseFound = "No database with name %v registered"
	//ErrNoDatabaseSelected is returned when no database was designated in the flag or env
	ErrNoDatabaseSelected = errors.New("No DB selected")
)

func init() {
	database = *flag.String("database", "mongodb", "Database to use")
}

//Init inits the selected DB in DefaultDb
func Init() error {
	if database == "" {
		return ErrNoDatabaseSelected
	}
	err := Set()
	if err != nil {
		return err
	}
	return DefaultDb.Init()
}

//Set the DefaultDb
func Set() error {
	if v, ok := DBTypes[database]; ok {
		DefaultDb = v
		return nil
	}
	return fmt.Errorf(ErrNoDatabaseFound, database)
}

//Register registers the database interface in the DBTypes
func Register(name string, db Database) {
	DBTypes[name] = db
}

//GetUserByName invokes DefaultDb method
func GetUserByName(n string) (m_user.User, error) {
	u, err := DefaultDb.GetUserByName(n)
	if err == nil {
	}
	return u, err
}

//GetUser invokes DefaultDb method
func GetUser(n string) (m_user.User, error) {
	u, err := DefaultDb.GetUser(n)
	if err == nil {
	}
	return u, err
}
