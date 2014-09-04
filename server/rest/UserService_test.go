package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/napping"
	"github.com/mikespook/gorbac"
	//"labix.org/v2/mgo/bson"
	"net/http/httptest"
	"testing"
)

func TestUserService (T *testing.T) {
	app := gin.Default()

	rbac = gorbac.New()
	rbac.Add("guest", []string{"create.user"}, nil)
	rbac.Set("master", []string{"list.user"}, []string{"guest"})

	db := NewDatabase("localhost:27017", "test")

	users := UserService{db:db}
	users.Bind(app)

	sessions := SessionService{db:db}
	sessions.Bind(app)

	db.Collection("Users", new(User))

	ts := httptest.NewServer(app)
	
	var user User
	var result_user User

	user.Username = "Travis"
	user.Public   = RandString()
	user.Password = Protect(user.Username + "password", user.Public)
	user.Role     = "master"

	res, err := napping.Post(ts.URL + "/users", &user, &result_user, nil)

	if err != nil {
		T.Fatal(err)
	}

	if res.Status() != 200 {
		T.Log("Expected status to be %s got %s", 200, res.Status())
		T.Fatal()
	}

	if len(result_user.Password) > 0 {
		T.Log("Password: %s", result_user.Password)
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

	if result_user.Username != user.Username {
		T.Fatal("Username does not match", result_user.Username , user.Username)
	}

	if result_user.Public != user.Public {
		T.Fatal("Public key does not match")
	}

	if len(result_user.ID) == 0  {
		T.Fatal("No id was returned")
	}

	var sess Session

	sess.Username = "Travis"

	res, err = napping.Post(ts.URL + "/session", &sess, &sess, nil)

	sess.Password = sess.HashPlainTextPassword("password")

	res, err = napping.Put(ts.URL + "/session/" +sess.Salt, &sess, &sess, nil)

	var user_list []User

	res, err = napping.Get(ts.URL + "/users", nil, &user_list, nil)

	s := napping.Session{}


	r := napping.Request{
		Method: "GET",
		Url:    ts.URL + "/users",
		Params: nil,
		Result: &user_list,
		Error:  nil,
	}

	r.Header.Add("session", sess.Salt)

	res, err = s.Send(&r)


	T.Log(user_list)

}
