package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/whyrusleeping/askeecs/server/kvstore"
	"bytes"
	"encoding/hex"
	"crypto/rand"
	"io"
	. "github.com/visionmedia/go-debug"
)

var debug = Debug("askeecs:session")

type SessionService struct {
	db *Database
}

type Session struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

func RandString() string {
	buf := new(bytes.Buffer)
	io.CopyN(buf, rand.Reader, 32)
	return hex.EncodeToString(buf.Bytes())
}

func (p *SessionService) Bind (app *gin.Engine) {
	app.POST("/session", p.CreateSession)
	app.DELETE("/session", p.DeleteSession)

	app.POST("/session/salt", p.CreateSessionSalt)
}

func (p *SessionService) CreateSession (c *gin.Context) {

	var sess Session

	if c.Bind(&sess) {
		salt := RandString();
		debug("Generated salt for %s [%s]", sess.Username, salt)
		kvstore.Set("SessionSalts", sess.Username, salt)

		sess.Salt = salt

		c.JSON(200, sess)
	}

}

func (p *SessionService) DeleteSession (c *gin.Context) {

}

func (p *SessionService) CreateSessionSalt(c *gin.Context) {
}
