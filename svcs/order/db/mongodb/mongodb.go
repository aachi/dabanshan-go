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
	name            string
	password        string
	host            string
	db              = "test"
	collections     = "orders"
	ErrInvalidHexID = errors.New("Invalid Id Hex")
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

// New Returns a new MongoOrder
func New() MongoOrder {
	u := m_order.New()
	return MongoOrder{
		Invoice: u,
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
	c := s.DB(db).C(collections)
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

// CreateUser Insert user to MongoDB
func (m *Mongo) CreateOrder(u *m_order.Invoice) (string, error) {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mu := New()
	mu.Invoice = *u
	mu.ID = id
	c := s.DB(db).C(collections)
	_, err := c.UpsertId(mu.ID, mu)
	if err != nil {
		return "", err
	}
	return mu.ID.Hex(), nil
}
