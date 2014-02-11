package main

import (
	"labix.org/v2/mgo/bson"
	"time"
	"bytes"
	"encoding/json"
	"io"
)

type User struct {
	_id bson.ObjectId
	Username string
	Password string `json:"-"`
	Salt string `json:"-"`
	login time.Time
}

func (u *User) New() I {
	return new(User)
}

func (u *User) GetID() bson.ObjectId {
	return u._id
}

func (u *User) JsonBytes() []byte {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(u)
	return buf.Bytes()
}

type AuthAttempt struct {
	Username string
	Password string
}

func AuthFromJson(r io.Reader) *AuthAttempt {
	a := new(AuthAttempt)
	dec := json.NewDecoder(r)
	err := dec.Decode(a)
	if err != nil {
		return nil
	}
	return a
}

