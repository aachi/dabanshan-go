package mongodb

import (
	"errors"
	"flag"
	"net/url"
	"time"

	m_order "github.com/laidingqing/dabanshan/svcs/order/model"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	name             string
	password         string
	host             string
	db               = "test"
	orderCollections = "orders"
	cartCollections  = "carts"
	ErrInvalidHexID  = errors.New("Invalid Id Hex")
)

func init() {
	name = *flag.String("mongouser", "", "Mongo user")
	password = *flag.String("mongopassword", "", "Mongo password")
	host = *flag.String("mongohost", "127.0.0.1:27017", "mongo host")
}

// Mongo meets the Database interface requirements
type Mongo struct {
	//Session is a MongoDB Session
	Session *mgo.Session
}

// MongoOrder is a wrapper for the users
type MongoOrder struct {
	m_order.Invoice `bson:",inline"`
	ID              bson.ObjectId `bson:"_id"`
}

// MongoCart is a wrapper for the users
type MongoCart struct {
	m_order.Cart `bson:",inline"`
	ID           bson.ObjectId `bson:"_id"`
}

// NewOrder Returns a new MongoOrder
func NewOrder() MongoOrder {
	u := m_order.New()
	return MongoOrder{
		Invoice: u,
	}
}

// NewCart ..
func NewCart() MongoCart {
	u := m_order.Cart{}
	return MongoCart{
		Cart: u,
	}
}

// Init MongoDB
func (m *Mongo) Init() error {
	u := getURL()
	var err error
	m.Session, err = mgo.DialWithTimeout(u.String(), time.Duration(5)*time.Second)
	if err != nil {
		return err
	}
	return m.EnsureIndexes()
}

// EnsureIndexes ensures userid is unique
func (m *Mongo) EnsureIndexes() error {
	s := m.Session.Copy()
	defer s.Close()
	i := mgo.Index{
		Key:        []string{"userId"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}
	c := s.DB(db).C(orderCollections)
	return c.EnsureIndex(i)
}

func getURL() url.URL {
	ur := url.URL{
		Scheme: "mongodb",
		Host:   host,
		Path:   db,
	}
	if name != "" {
		u := url.UserPassword(name, password)
		ur.User = u
	}
	return ur
}

// CreateOrder Insert user to MongoDB
func (m *Mongo) CreateOrder(u *m_order.Invoice) (string, error) {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mu := NewOrder()
	mu.Invoice = *u
	mu.ID = id
	c := s.DB(db).C(orderCollections)
	_, err := c.UpsertId(mu.ID, mu)
	if err != nil {
		return "", err
	}
	return mu.ID.Hex(), nil
}

// GetOrders 根据用户查询订单列表.
func (m *Mongo) GetOrders(usrID string) ([]m_order.Invoice, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB(db).C(orderCollections)
	var orders []m_order.Invoice
	err := c.Find(bson.M{"userId": usrID}).All(&orders)

	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrder 根据用户查询订单.
func (m *Mongo) GetOrder(id string) (m_order.Invoice, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB(db).C(orderCollections)
	var order m_order.Invoice
	err := c.FindId(id).One(&order)

	if err != nil {
		return m_order.Invoice{}, err
	}
	return order, nil
}

// GetCartItems ..
func (m *Mongo) GetCartItems(userID string) ([]m_order.Cart, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB(db).C(cartCollections)
	var cartItems []m_order.Cart
	err := c.Find(bson.M{"userID": userID}).All(&cartItems)

	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

// AddCart ..
func (m *Mongo) AddCart(cart *m_order.Cart) (string, error) {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mu := NewCart()
	mu.Cart = *cart
	mu.ID = id
	c := s.DB(db).C(cartCollections)
	_, err := c.UpsertId(mu.ID, mu)
	if err != nil {
		return "", err
	}
	return mu.ID.Hex(), nil
}

// RemoveCartItem ..
func (m *Mongo) RemoveCartItem(cartID string) (bool, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB(db).C(cartCollections)
	err := c.RemoveId(cartID)
	if err != nil {
		return false, err
	}
	return true, nil
}
