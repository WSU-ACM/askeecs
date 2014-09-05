package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/whyrusleeping/askeecs/server/kvstore"
	"bytes"
	"labix.org/v2/mgo/bson"
	"encoding/hex"
	"crypto/sha256"
	"crypto/rand"
	"encoding/json"
	"io"
	. "github.com/visionmedia/go-debug"
)

var session_debug = Debug("askeecs:session")

type SessionService struct {
	db *Database
}

type Session struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Public   string `json:"public"`
	Valid    bool   `json:"valid"`
	Role     string `json:"role"`
}

func (s *Session) ValidUser (user *User) bool {
	return s.ValidateHash(user.Password)
}

func (s *Session) ValidateHash (password string) bool {
	return s.Password == Protect(password, s.Salt)
}

func (s *Session) ValidatePlain (password string) bool {
	return s.Password == s.HashPlainTextPassword(password)
}

func (s *Session) HashPlainTextPassword (password string) string {
	passhash := Protect(s.Username + password, s.Public)
	passhash  = Protect(passhash, s.Salt)

	return passhash
}

func GetSession (c *gin.Context) (Session, bool) {
	var sess Session
	session_id := c.Request.Header.Get("session")

	I, found := kvstore.Get("Session", session_id)

	if !found {
		c.JSON(501, gin.H{"message": "Session does not exist"})
		return sess, false
	}

	sess = I.(Session)

	return sess, true
}


func (s *Session) Marshal() []byte {
	b, err := json.Marshal(s)

	if err != nil {
		session_debug("Error Marshalling")
	}

	return b
}

func (s *Session) Decode(b []byte) {
	dec := json.NewDecoder(bytes.NewBuffer(b))

	if err := dec.Decode(s); err == nil {
		session_debug("Worked")
	} else if err != nil {
		session_debug("Error")
		panic(err)
		return
	}
}

func Protect(pass, salt string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

func RandString() string {
	buf := new(bytes.Buffer)
	io.CopyN(buf, rand.Reader, 32)
	return hex.EncodeToString(buf.Bytes())
}

func (p *SessionService) Bind (app *gin.Engine) {
	app.POST   ("/session",     p.CreateSession)
	app.DELETE ("/session/:id", p.DeleteSession)
	app.PUT    ("/session/:id", p.ValidateSession)
}

func (p *SessionService) CreateSession (c *gin.Context) {

	var sess Session

	if c.Bind(&sess) {
		result := p.db.collections["Users"].FindOneWhere(bson.M{"username": sess.Username})
		
		salt     := RandString()
		sess.Salt = salt

		session_debug("Generated salt for %s [%s]", sess.Username, salt)

		kvstore.Set("Session", sess.Username+":salt", salt)
		kvstore.Set("Session", salt, nil)

		if result == nil {
			sess.Public= RandString()
		} else {
			user := result.(*User)
			sess.Public = user.Public
		}

		kvstore.Set("Session", sess.Username+":valid_user", result != nil)

		c.JSON(200, sess)
	}

}

func (p *SessionService) ValidateSession(c *gin.Context) {
	var sess Session

	if c.Bind(&sess) {
		valid_session, found := kvstore.Get("Session", sess.Username+":valid_user")

		if valid_session == false || found == false {
			sess.Valid = false
			c.JSON(200, sess)
			return
		}

		result := p.db.collections["Users"].FindOneWhere(bson.M{"username": sess.Username})
		
		if result == nil {
			c.JSON(500, gin.H{"message": "Could not find the matching user"})
			return
		}

		user := result.(*User)

		salt, found := kvstore.Get("Session", sess.Username+":salt")

		if !found {
			c.JSON(500, gin.H{"message": "Could not find salt"})
			return
		}

		if sess.ValidUser(user) {
			sess.Valid = true
			sess.Role  = user.Role
			kvstore.Set("Session", salt.(string), sess)
		} else {
			sess.Valid = false
			kvstore.Set("Session", salt.(string), nil)
		}

		c.JSON(200, sess)

	}

}

func (p *SessionService) DeleteSession (c *gin.Context) {
	c.JSON(200, gin.H{"status":"ok"})
}
