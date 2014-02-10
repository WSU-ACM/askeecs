package main

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type User struct {
	_id bson.ObjectId
	Username string
	Password string
	Salt string
	login time.Time
}

func (u *User) New() I {
	return new(User)
}

func (u *User) GetID() bson.ObjectId {
	return u._id
}

