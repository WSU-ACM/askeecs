package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/napping"
	//"labix.org/v2/mgo/bson"
	"net/http/httptest"
	"testing"
)


func TestUserService (T *testing.T) {
	app := gin.Default()

	db := NewDatabase("localhost:27017", "test")

	users := UserService{db:db}
	users.Bind(app)

	db.Collection("Users", new(User))
	db.db.C("Users").DropCollection()

	ts := httptest.NewServer(app)
	
	var user User
	var result_user User

	user.Username = "Travis"
	user.Public   = RandString()
	user.Password = Protect(user.Username + "password", user.Public)

	res, err := napping.Post(ts.URL + "/users", &user, &result_user, nil)

	if err != nil {
		T.Fatal(err)
	}

	if res.Status() != 200 {
		T.Log("Expected status to be %s got %s", 200, res.Status())
		T.Fatal()
	}

	if len(result_user.Password) > 0 {
		T.Fatal("Password was present in the response")
	}

	res, err = napping.Get(ts.URL + "/users/" + result_user.ID.Hex(), nil, &result_user, nil)

	if res.Status() != 200 {
		T.Log("Expected status to be %s got %s", 200, res.Status())
		T.Fatal()
	}

	if len(result_user.Password) > 0 {
		T.Fatal("Password was present in the response")
	}

}
