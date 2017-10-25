package mongodb

import (
	"errors"
	"flag"
	"net/url"
	"time"

	m_product "github.com/laidingqing/dabanshan/svcs/product/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	name            string
	password        string
	host            string
	db              = "test"
	ErrInvalidHexID = errors.New("Invalid Id Hex")
)

func init() {
	name = *flag.String("mongouser", "", "Mongo user")
	password = *flag.String("mongopassword", "", "Mongo password")
	host = *flag.String("mongohost", "127.0.0.1:27017", "mongo host")
}

// Mongo ...
type Mongo struct {
	//Session is a MongoDB Session
	Session *mgo.Session
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

// CreateProduct ...
func (m *Mongo) CreateProduct(p *m_product.Product) error {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mp := New()
	mp.Product = *p
	mp.ID = id
	c := s.DB(db).C("products")
	_, err := c.UpsertId(mp.ID, mp)
	if err != nil {
		return err
	}
	mp.Product.ProductID = mp.ID.Hex()
	*p = mp.Product
	return nil
}

// MongoProduct is a wrapper for the users
type MongoProduct struct {
	m_product.Product `bson:",inline"`
	ID                bson.ObjectId `bson:"_id"`
}

// New Returns a new MongoProduct
func New() MongoProduct {
	p := m_product.New()
	return MongoProduct{
		Product: p,
	}
}

// EnsureIndexes ensures userid is unique
func (m *Mongo) EnsureIndexes() error {
	s := m.Session.Copy()
	defer s.Close()
	i := mgo.Index{
		Key:        []string{"userid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}
	c := s.DB("").C("products")
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
