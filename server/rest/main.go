package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/mikespook/gorbac"
)

type Api struct {
	// API info
	ConnString string
	DBName string
}


var rbac *gorbac.Rbac


func (a *Api) Middleware (c *gin.Context) {

}

func (a *Api) Run (binding string) {
	app := gin.Default()
	app.Use(a.Middleware)

	rbac = gorbac.New()

	// db info
	db := NewDatabase(a.ConnString, a.DBName)

	sessions := SessionService{db:db}
	sessions.Bind(app)

	users := UserService{db:db}
	users.Bind(app)

	rbac.Add("guest", []string{"create.user"}, nil)
	rbac.Set("master", []string{"list.user"}, []string{"guest"})

	app.Run(binding)
}


