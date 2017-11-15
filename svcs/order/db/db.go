package db

import (
	"errors"
	"fmt"
	corelog "log"

	m_order "github.com/laidingqing/dabanshan/svcs/order/model"
	"github.com/laidingqing/dabanshan/utils"
)

// Database represents a simple interface so we can switch to a new system easily
type Database interface {
	Init() error
	CreateOrder(*m_order.Invoice) (string, error)
	GetOrdersByUser(usrID string, page utils.Pagination) (utils.Pagination, error)
	GetOrdersByTenant(usrID string, page utils.Pagination) (utils.Pagination, error)
	GetOrder(id string) (m_order.Invoice, error)
	AddCart(cart *m_order.Cart) (string, error)
	RemoveCartItem(cartID string) (bool, error)
	GetCartItems(userID string) ([]m_order.Cart, error)
	UpdateQuantity(cart *m_order.Cart) (m_order.Cart, error)
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
	//database = *flag.String("database", "mongodb", "Database to use")
	database = "mongodb"
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

// GetOrdersByUser ...
func GetOrdersByUser(usrID string, page utils.Pagination) (utils.Pagination, error) {
	return DefaultDb.GetOrdersByUser(usrID, page)
}

// GetOrdersByTenant ...
func GetOrdersByTenant(tenantID string, page utils.Pagination) (utils.Pagination, error) {
	return DefaultDb.GetOrdersByTenant(tenantID, page)
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
	corelog.Print("cartID is " + cartID)
	return DefaultDb.RemoveCartItem(cartID)
}

// GetCartItems ..
func GetCartItems(userID string) ([]m_order.Cart, error) {
	return DefaultDb.GetCartItems(userID)
}

// UpdateQuantity ..
func UpdateQuantity(cart *m_order.Cart) (m_order.Cart, error) {
	return DefaultDb.UpdateQuantity(cart)
}
