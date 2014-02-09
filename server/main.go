package main

import (
	"time"
	"github.com/codegangsta/martini"
	"github.com/whyrusleeping/jadb"
	"encoding/json"
)

type Question struct {
	ID string
	Title string
	Author string
	Tags []string
	Score int
	Timestamp time.Time
	Body string
	Responses []*Response
}

func (q *Question) New() jadb.I {
	return new(Question)
}

func (q *Question) GetID() string {
	return q.ID
}

type Response struct {
	Author string
	Timestamp time.Time
	Score int
	Body string
}

var db *jadb.Jadb

func main() {
	db = jadb.NewJadb("data")
	questions := db.Collection("Questions", new(Question))
	m := martini.Classic()
	m.Get("/q/:id", func(params martini.Params) (int, string) {
		id := params["id"]
		q,ok := questions.FindByID(id).(*Question)
		if !ok || q == nil {
			return 404,""
		}
		b,_ := json.Marshal(q)
		return 200, string(b)
	})
	m.Run()
}
