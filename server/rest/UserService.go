
package rest

import (
	"github.com/gin-gonic/gin"
//	"github.com/whyrusleeping/askeecs/server/kvstore"
	"labix.org/v2/mgo/bson"
)

type UserService struct {
	db *Database
}

type User struct {
	ID bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty" bson: "password"`
	Public   string `json:"public"`
}

func (this *User) GetID() bson.ObjectId {
	return this.ID
}

func (this *User) Sanitize() User {
	var user User

	user = *this
	user.Password = ""

	return user
}

func (this *User) New() I {
	return new(User)
}

func (p *UserService) Bind (app *gin.Engine) {
	p.db.Collection("Users", new(User))
	app.GET("/users", p.ListUsers)
	app.GET("/users/:id", p.GetUser)
	app.POST("/users", p.CreateUser)
}

func (p *UserService) ListUsers (c *gin.Context) {
	list := p.db.collections["Users"].FindWhere(bson.M{})
	if list == nil {
		c.JSON(404, gin.H{"message": "no records found"})
		return
	}

	c.JSON(200, list)
}

func (p *UserService) GetUser(c *gin.Context) {
	var user_id = c.Params.ByName("id")

	result := p.db.collections["Users"].FindByID(bson.ObjectIdHex(user_id))

	if result == nil {
		c.JSON(500, gin.H{"message": "Could not find user"})
	} else {
		user := result.(*User)
		c.JSON(200, user.Sanitize())
	}
}

func (p *UserService) CreateUser(c *gin.Context) {
	var user User
	var err error

	if c.Bind(&user) {
		user.ID = bson.NewObjectId()
		db_debug("%s", user)
		err = p.db.collections["Users"].Save(&user)

		if err != nil {
			c.JSON(500, gin.H{"message": "error making user"})
			panic(err)
			return
		}

		c.JSON(200, user.Sanitize())

	}
}


