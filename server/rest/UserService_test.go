package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/napping"
	"github.com/mikespook/gorbac"
	//"labix.org/v2/mgo/bson"
	"net/http/httptest"
//	"net/http"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserService (T *testing.T) {
	app := gin.New()

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

	Convey("Create new user", T, func() {

		res, _ := napping.Post(ts.URL + "/users", &user, &result_user, nil)

		Convey("User should be created", func() {
			So(res.Status(), ShouldEqual, 200)
		})

		Convey("Password should not be sent back with the response", func() {
			So(result_user.Password, ShouldBeBlank)
		})

	})


	Convey("User should be retrieveable with their ID", T, func() {
		res, _ := napping.Get(ts.URL + "/users/" + result_user.ID.Hex(), nil, &result_user, nil)

		Convey("User should exist", func() {
			So(res.Status(), ShouldEqual, 200)
		})

		Convey("Password should not be sent back with the response", func() {
			So(result_user.Password, ShouldBeBlank)
		})

		Convey("The username should be the same", func() {
			So(user.Username, ShouldEqual, result_user.Username)
		})

		Convey("The public key should be the same", func() {
			So(user.Public, ShouldEqual, result_user.Public)
		})

		Convey("User should get an ID", func() {
			So(result_user.ID.Hex() , ShouldNotBeBlank)
		})
	})

	var sess Session
	sess.Username = "Travis"

	Convey("User should be able to login", T, func() {
		res, _ := napping.Post(ts.URL + "/session", &sess, &sess, nil)

		Convey("A new session should be created", func() {
			So(res.Status(), ShouldEqual, 200)
		})

		Convey("A salt should be returned", func() {
			So(sess.Salt, ShouldNotBeBlank)
		})

		Convey("A public key should be returned", func() {
			So(sess.Public, ShouldNotBeBlank)
		})

		sess.Password = sess.HashPlainTextPassword("password")

		res, _ = napping.Put(ts.URL + "/session/" +sess.Salt, &sess, &sess, nil)

		Convey("Session should be created", func() {
			So(res.Status(), ShouldEqual, 200)
		})

		Convey("Session should be validated", func() {
			So(sess.Valid, ShouldBeTrue)
		})
/*
		var user_list []User

		res, _ = napping.Get(ts.URL + "/users", nil, &user_list, nil)

		if res.Status() != 501 {
			T.Logf("Exected status to be %s got %s", 501, res.Status())
			T.Fatal()
		}

		s := napping.Session{}

		r := napping.Request{
			Method: "GET",
			Url:    ts.URL + "/users",
			Params: nil,
			Result: &user_list,
			Error:  nil,
		}

		r.Header = &http.Header{}

		r.Header.Add("session", sess.Salt)

		res, _ = s.Send(&r)

		if res.Status() != 200 {
			T.Logf("Exected status to be %s got %s", 200, res.Status())
		}
		*/
		
	})

}
