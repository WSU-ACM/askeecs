package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/whyrusleeping/askeecs/server/kvstore"
	"strconv"
)

type PingPong struct {
	db *Database
}

func (p *PingPong) Bind (app *gin.Engine) {
	app.GET("/ping", p.Get)
}

func (p *PingPong) Get (c *gin.Context) {
	v, found := kvstore.Get("tmp", "count")

	i,_ := v.(int)

	if !found {
		kvstore.Set("tmp", "count", 1)
	} else {
		kvstore.Set("tmp", "count", i + 1)
	}

	c.String(200, "pong" + strconv.Itoa(i))
}


