package mongodb

import (
	"errors"
	"flag"
	"net/url"
	"time"

	m_user "github.com/laidingqing/dabanshan/svcs/user/model"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	name            string
	password        string
	host            string
	db              = "test"
	collections     = "users"
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

// MongoUser is a wrapper for the users
type MongoUser struct {
	m_user.User `bson:",inline"`
	ID          bson.ObjectId `bson:"_id"`
}

// New Returns a new MongoUser
func New() MongoUser {
	u := m_user.New()
	return MongoUser{
		User: u,
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
		Key:        []string{"username"},
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

// GetUserByName Get user by their name
func (m *Mongo) GetUserByName(name string) (m_user.User, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB(db).C(collections)
	mu := New()
	err := c.Find(bson.M{"username": name}).One(&mu)
	mu.UserID = mu.ID.Hex()
	return mu.User, err
}

// CreateUser Insert user to MongoDB
func (m *Mongo) CreateUser(u *m_user.User) (string, error) {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mu := New()
	mu.User = *u
	mu.ID = id
	c := s.DB(db).C(collections)
	_, err := c.UpsertId(mu.ID, mu)
	if err != nil {
		return "", err
	}
	return mu.ID.Hex(), nil
}

// GetUser Get user by their object id
func (m *Mongo) GetUser(id string) (m_user.User, error) {
	s := m.Session.Copy()
	defer s.Close()
	if !bson.IsObjectIdHex(id) {
		return m_user.New(), errors.New("Invalid Id Hex")
	}
	c := s.DB(db).C("users")
	mu := New()
	err := c.FindId(bson.ObjectIdHex(id)).One(&mu)
	mu.UserID = mu.ID.Hex()
	return mu.User, err
}
