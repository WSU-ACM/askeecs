package main

import (
	"time"
	"fmt"
	"net/http"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo/bson"
	"encoding/json"
	"encoding/hex"
	"log"
	"crypto/rand"
	"bytes"
	"io"
	"io/ioutil"
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

type AEServer struct {
	db *Database
	questions *Collection
	users *Collection

	tokens map[string]*User
	salts map[string]string
}

func NewServer() *AEServer {
	s := new(AEServer)
	s.db = NewDatabase("localhost:27017")
	s.questions = s.db.Collection("Questions", new(Question))
	s.users = s.db.Collection("Users", new(User))
	s.tokens = make(map[string]*User)
	return s
}

func (s *AEServer) GetSessionToken() string {
	buf := new(bytes.Buffer)
	io.CopyN(buf, rand.Reader, 32)
	return hex.EncodeToString(buf.Bytes())
}

func (s *AEServer) HandlePostQuestion(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	//Verify user account or something
	login := session.Get("Login")
	if login == nil {
		log.Printf("Not logged in!!")
		w.WriteHeader(404)
		return
	}
	tok := login.(string)
	user, ok := s.tokens[tok]
	if !ok {
		log.Printf("Invalid cookie!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	var q Question
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&q)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}
	q.ID = bson.NewObjectId()
	q.Author = user.Username
	fmt.Println(q.ID)
	err = s.questions.Save(&q)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusTeapot)
	}
	str := hex.EncodeToString([]byte(q.ID))
	w.Write([]byte(str))
}

func (s *AEServer) HandleGetQuestions()(int,string) {
	q := s.questions.FindWhere(bson.M{})
	if q == nil {
		return 404,""
	}
	b,_ := json.Marshal(q)
	return 200, string(b)
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

func (s *AEServer) HandleLogout(session sessions.Session) {
	toki := session.Get("Login")
	if toki == nil {
		return
	}
	tok := toki.(string)
	delete(s.tokens, tok)
	session.Delete("Login")
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
		time.Sleep(time.Second)
		return 404, "Failed"
	}
	users := s.users.FindWhere(bson.M{"username":a.Username})
	if len(users) == 0 {
		fmt.Println("User not found.")
		time.Sleep(time.Second)
		return http.StatusUnauthorized, "Invalid Username or Password."
	}

	user, _ := users[0].(*User)

	if user.Password != a.Password {
		fmt.Println("Invalid password.")
		time.Sleep(time.Second)
		return http.StatusUnauthorized, "Invalid Username or Password."
	}

	tok := s.GetSessionToken()
	user.login = time.Now()
	for _,ok := s.tokens[tok]; ok; tok = s.GetSessionToken() {}
	s.tokens[tok] = user

	session.Set("Login", tok);

	fmt.Println("Logged in!");
	return 200, "OK"
}

func (s *AEServer) HandleMe(session sessions.Session) (int, string) {
	return 200, "Nothing here"
}

func (s *AEServer) HandleRegister(r *http.Request) (int, string) {
	var a AuthAttempt
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&a)
	if err != nil {
		fmt.Println(err)
		return 404, "Register Failed"
	}

	user := new(User)
	user.Password = a.Password
	user.Username = a.Username
	user._id = bson.NewObjectId()

	s.users.Save(user)
	return 200,"Success!"
}

func main() {
	s := NewServer()
	m := martini.Classic()


	secret,err := ioutil.ReadFile(".secret")
	if err != nil {
		panic(err)
	}
	store := sessions.NewCookieStore(secret)
    m.Use(sessions.Sessions("ask_eecs_auth_session", store))

	m.Get("/q/:id", s.HandleGetQuestion)
	m.Post("/q", s.HandlePostQuestion)
	m.Get("/q", s.HandleGetQuestions)
	m.Post("/login", s.HandleLogin)
	m.Post("/register", s.HandleRegister)
	m.Post("/logout", s.HandleLogout)
	m.Post("/me", s.HandleMe);
	m.Run()
}
