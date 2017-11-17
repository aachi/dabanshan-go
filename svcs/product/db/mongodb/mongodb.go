package mongodb

import (
	"flag"
	"net/url"
	"strconv"
	"time"

	m_product "github.com/laidingqing/dabanshan/svcs/product/model"
	"github.com/laidingqing/dabanshan/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	name        string
	password    string
	host        string
	db          = "test"
	collections = "products"
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

// MongoProduct is a wrapper for the users
type MongoProduct struct {
	m_product.Product `bson:",inline"`
	ID                bson.ObjectId `bson:"_id"`
}

// NewOrder Returns a new MongoOrder
func NewProduct() MongoProduct {
	p := m_product.New()
	return MongoProduct{
		Product: p,
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

// CreateProduct ...
func (m *Mongo) CreateProduct(p *m_product.Product) (string, error) {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mp := NewProduct()
	mp.Product = *p
	mp.ID = id
	c := s.DB(db).C(collections)
	_, err := c.UpsertId(mp.ID, mp)
	if err != nil {
		return "", err
	}
	mp.Product.ID = mp.ID.Hex()
	*p = mp.Product
	return mp.ID.Hex(), nil
}

// UploadGfs ...
func (m *Mongo) UploadGfs(body []byte) (string, error) {
	gf, _ := utils.NewGlowFlake(1, 1)
	id, _ := gf.NextId()
	fsid := strconv.FormatInt(id, 10)
	s := m.Session.Copy()
	defer s.Close()
	fs, err := s.DB(db).GridFS("fs").Create(fsid)
	if err != nil {
		return "", err
	}
	defer fs.Close()
	fs.SetName(fsid)
	if _, err := fs.Write(body); err != nil {
		return "", err
	}
	return fsid, nil
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
