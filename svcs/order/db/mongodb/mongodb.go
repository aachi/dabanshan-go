package mongodb

import (
	"errors"
	"flag"
	corelog "log"
	"net/url"
	"os"
	"time"

	"github.com/go-kit/kit/log"
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

var logger log.Logger

func init() {
	name = *flag.String("mongouser", "", "Mongo user")
	password = *flag.String("mongopassword", "", "Mongo password")
	host = *flag.String("mongohost", "127.0.0.1:27017", "mongo host")

	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
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
		Unique:     false,
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
	u.CreatedAt = time.Now()
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
	// not debug data.
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
	mu.ID = id
	mu.Cart = *cart
	c := s.DB(db).C(cartCollections)
	_, err := c.UpsertId(id, mu)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

// RemoveCartItem ..
func (m *Mongo) RemoveCartItem(cartID string) (bool, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB(db).C(cartCollections)
	corelog.Print("cartID:" + cartID)
	err := c.RemoveId(bson.ObjectIdHex(cartID))
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateQuantity update quantity of cartitem
func (m *Mongo) UpdateQuantity(cart *m_order.Cart) (m_order.Cart, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB(db).C(cartCollections)

	corelog.Print("cartID:" + cart.CartID)
	var item m_order.Cart
	err := c.FindId(bson.ObjectIdHex(cart.CartID)).One(&item)
	if err != nil {
		return m_order.Cart{}, err
	}
	item.CartID = cart.CartID
	item.Quantity = cart.Quantity
	err = c.UpdateId(bson.ObjectIdHex(item.CartID), item)
	if err != nil {
		return m_order.Cart{}, err
	}
	return item, nil
}
