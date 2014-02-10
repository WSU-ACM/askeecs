package main

import (
	"time"
	"fmt"
	"net/http"
	"github.com/codegangsta/martini"
	"labix.org/v2/mgo/bson"
	"encoding/json"
	"encoding/hex"
	"log"
)

type Score struct {
	Up int
	Down int
}

type Question struct {
	ID bson.ObjectId "_id,omitempty"
	Title string
	Author string
	Tags []string
	Score Score
	Timestamp time.Time
	Body string
	Responses []*Response
	Comments []*Comment
}

func (q *Question) New() I {
	return new(Question)
}

func (q *Question) GetID() bson.ObjectId {
	return q.ID
}

type Response struct {
	ID bson.ObjectId
	Author string
	Timestamp time.Time
	Score Score
	Body string
	Comments []*Comment
}

type Comment struct {
	ID bson.ObjectId
	Timestamp time.Time
	Author string
	Content string
	Score Score
}

type AEServer struct {
	db *Database
	questions *Collection
}

func NewServer() *AEServer {
	s := new(AEServer)
	s.db = NewDatabase("localhost:27017")
	s.questions = s.db.Collection("Questions", new(Question))
	return s
}

func (s *AEServer) HandlePostQuestion(w http.ResponseWriter, r *http.Request) {
	//Verify user account or something
	var q Question
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&q)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}
	q.ID = bson.NewObjectId()
	fmt.Println(q.ID)
	err = s.questions.Save(&q)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusTeapot)
	}
	str := hex.EncodeToString([]byte(q.ID))
	w.Write([]byte(str))
}

func (s *AEServer) HandleGetQuestion(params martini.Params) (int,string) {
	id := params["id"]
	hid := bson.ObjectIdHex(id)
	fmt.Println(hid)
	q,ok := s.questions.FindByID(hid).(*Question)
	if !ok || q == nil {
		return 404,""
	}
	b,_ := json.Marshal(q)
	return 200, string(b)
}

func main() {
	s := NewServer()
	m := martini.Classic()
	m.Get("/q/:id", s.HandleGetQuestion)
	m.Post("/q", s.HandlePostQuestion)
	m.Run()
}
