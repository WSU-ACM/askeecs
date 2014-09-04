package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/napping"
	"labix.org/v2/mgo/bson"
	"net/http/httptest"
	"testing"
)


func TestSessionService (T *testing.T) {
	app := gin.Default()

	db := NewDatabase("localhost:27017", "test")

	sessions := SessionService{db:db}
	sessions.Bind(app)

	db.Collection("Users", new(User))
	db.db.C("Users").DropCollection()

	ts := httptest.NewServer(app)
	
	var user User

	user.ID       = bson.NewObjectId()
	user.Username = "Travis"
	user.Public   = RandString()
	user.Password = Protect(user.Username + "password", user.Public)
	user.Role     = "master"
	
	db.collections["Users"].Save(&user)

	//

	var sess Session

	sess.Username = "Travis"

	res, err := napping.Post(ts.URL + "/session", &sess, &sess, nil)

	if err != nil {
		T.Fatal(err)
	}

	if res.Status() != 200 {
		T.Log("Expected status to be %s got %s", 200, res.Status())
		T.Fatal()
	}

	if user.Public != sess.Public {
		T.Log(user.Public)
		T.Log(sess.Public)
		T.Fatal("Public keys do not match")
	}

	if user.Username != sess.Username {
		T.Fatal("Usernames do not match")
	}
	
	if len(sess.Password) > 0 {
		T.Fatal("Password was present")
	}

	if len(sess.Salt) != 256/4 {
		T.Log("Expected a length of %v got %v", 256/4, len(sess.Salt))
		T.Fatal("Salt is not valid")
	}

	sess.Password = sess.HashPlainTextPassword("password")

	res, err = napping.Put(ts.URL + "/session/" +sess.Salt, &sess, &sess, nil)

	if res.Status() != 200 {
		T.Log("Expected status to be %s got %s", 200, res.Status())
		T.Fatal()
	}

	if sess.Valid == false {
		T.Log("Failed to validate session")
		T.Fatal()
	}

	res, err = napping.Delete(ts.URL + "/session/" + sess.Salt, &sess, nil)

	if res.Status() != 200 {
		T.Log("Expected status to be %s got %s", 200, res.Status())
		T.Fatal()
	}
}
