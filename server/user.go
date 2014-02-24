package main

import (
	"labix.org/v2/mgo/bson"
	"time"
	"encoding/json"
	"io"
)

type User struct {
	ID bson.ObjectId "_id,omitempty"
	Username string
	Password string `json:"-"`
	Salt string `json:"-"`
	login time.Time
}

func (u *User) New() I {
	return new(User)
}

func (u *User) GetID() bson.ObjectId {
	return u.ID
}

type AuthAttempt struct {
	Username string
	Password string
	Salt string
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

func (u *User) MakeComment(r io.Reader) *Comment {
	c := CommentFromJson(r)
	if c == nil {
		return nil
	}
	c.Author = u.Username
	c.ID = bson.NewObjectId()
	c.Timestamp = time.Now()
	return c
}

func (u *User) MakeRespose(r io.Reader) *Response {
	resp := ResponseFromJson(r)
	if r == nil {
		return nil
	}

	resp.Author = u.Username
	resp.ID = bson.NewObjectId()
	resp.Timestamp = time.Now()
	return resp
}
