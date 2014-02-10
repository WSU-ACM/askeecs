package main

import (
	"time"
	"fmt"
	"net/http"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo/bson"
	"encoding/json"
	//"encoding/hex"
)

type Score struct {
	Up int
	Down int
}

type Question struct {
	ID bson.ObjectId
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
	s.questions.Save(&q)
	w.Write([]byte(q.ID))
}

func (s *AEServer) HandleGetQuestion(params martini.Params) (int,string) {
	id := params["id"]
	q,ok := s.questions.FindByID(bson.ObjectIdHex(id)).(*Question)
	if !ok || q == nil {
		return 404,""
	}
	b,_ := json.Marshal(q)
	return 200, string(b)
}

type AuthAttempt struct {
	Username string
	Password string
}


func (s *AEServer) HandleLogin(r *http.Request, params martini.Params, session sessions.Session) (int,string) {
	var a AuthAttempt
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&a)
	if err != nil {
		fmt.Println(err)
		return 404, "Failed"
	}

	fmt.Println(a.Username);
	fmt.Println(a.Password);

	session.Set("Login", "1");

	return 200, "OK"

}

func (s *AEServer) HandleMe(session sessions.Session) (int,string) {
	fmt.Println(session.Get("Login"))
	return 200, "OK"
}

func main() {
	s := NewServer()
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("secret123"))
    m.Use(sessions.Sessions("my_session", store))
	m.Get("/q/:id", s.HandleGetQuestion)
	m.Post("/q", s.HandlePostQuestion)
	m.Post("/login", s.HandleLogin);
	m.Post("/me", s.HandleMe);
	m.Run()
}
