package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/whyrusleeping/askeecs/server/kvstore"
	"strconv"
)

type BeepBoop struct {
	db *Database
}

func (p *BeepBoop) Bind (app *gin.Engine) {
	app.GET("/beep", p.Get)
}

func (p *BeepBoop) Get (c *gin.Context) {
	v, found := kvstore.Get("tmp", "count")

	i,_ := v.(int)

	if !found {
		kvstore.Set("tmp", "count", 1)
	} else {
		kvstore.Set("tmp", "count", i + 1)
	}

	c.String(200, "boop" + strconv.Itoa(i))
}


