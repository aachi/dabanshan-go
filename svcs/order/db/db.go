package db

import (
	"errors"
	"flag"
	"fmt"

	m_order "github.com/laidingqing/dabanshan/svcs/order/model"
)

// Database represents a simple interface so we can switch to a new system easily
type Database interface {
	Init() error
	CreateOrder(*m_order.Invoice) (string, error)
	GetOrders(usrID string) ([]m_order.Invoice, error)
	GetOrder(id string) (m_order.Invoice, error)
	AddCart(cart *m_order.Cart) (string, error)
	RemoveCartItem(cartID string) (bool, error)
	GetCartItems(userID string) ([]m_order.Cart, error)
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

// CreateOrder db operator
func CreateOrder(mo *m_order.Invoice) (string, error) {
	return DefaultDb.CreateOrder(mo)
}

// GetOrders ...
func GetOrders(usrID string) ([]m_order.Invoice, error) {
	return DefaultDb.GetOrders(usrID)
}

// GetOrder ...
func GetOrder(id string) (m_order.Invoice, error) {
	return DefaultDb.GetOrder(id)
}

// AddCart ..
func AddCart(cart *m_order.Cart) (string, error) {
	return DefaultDb.AddCart(cart)
}

// RemoveCartItem ..
func RemoveCartItem(cartID string) (bool, error) {
	return DefaultDb.RemoveCartItem(cartID)
}

// GetCartItems ..
func GetCartItems(userID string) ([]m_order.Cart, error) {
	return DefaultDb.GetCartItems(userID)
}
