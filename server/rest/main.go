package rest

import (
	"github.com/gin-gonic/gin"
)

type Api struct {
	// API info
	ConnString string
	DBName string
}

func (a *Api) Run (binding string) {
	app := gin.Default()

	// db info
	db := NewDatabase(a.ConnString, a.DBName)

	sessions := SessionService{db:db}
	sessions.Bind(app)

	users := UserService{db:db}
	users.Bind(app)

	app.Run(binding)
}


