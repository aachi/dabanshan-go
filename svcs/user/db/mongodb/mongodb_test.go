package mongodb

import (
	"os"
	"testing"

	users "github.com/laidingqing/dabanshan/svcs/user/model"
	"gopkg.in/mgo.v2/dbtest"
)

var (
	TestMongo  = Mongo{}
	TestServer = dbtest.DBServer{}
	TestUser   = users.User{
		FirstName: "firstname",
		LastName:  "lastname",
		Username:  "username",
		Password:  "blahblah",
	}
)

func init() {
	TestServer.SetPath("/tmp")
}

func TestMain(m *testing.M) {
	TestMongo.Session = TestServer.Session()
	TestMongo.EnsureIndexes()
	TestMongo.Session.Close()
	exitTest(m.Run())
}

func exitTest(i int) {
	TestServer.Wipe()
	TestServer.Stop()
	os.Exit(i)
}

func TestInit(t *testing.T) {
	err := TestMongo.Init()
	if err.Error() != "no reachable servers" {
		t.Error("expecting no reachable servers error")
	}
}

func TestGetUser(t *testing.T) {
	TestMongo.Session = TestServer.Session()
	defer TestMongo.Session.Close()
	_, err := TestMongo.GetUser(TestUser.UserID)
	if err != nil {
		t.Error(err)
	}
}
