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
)

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
	q.Timestamp = time.Now()

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
		return 404, "{\"Message\":\"Login Failed\"}"
	}
	users := s.users.FindWhere(bson.M{"username":a.Username})
	if len(users) == 0 {
		fmt.Println("User not found.")
		time.Sleep(time.Second)
		return 401, "{\"Message\":\"Invalid Username or Password\"}"
	}

	user, _ := users[0].(*User)

	fmt.Println(user.Password)
	if user.Password != a.Password {
		fmt.Println("Invalid password.")
		time.Sleep(time.Second)
		return http.StatusUnauthorized, "{\"Message\":\"Invalid Username or Password.\"}"
	}

	tok := s.GetSessionToken()
	user.login = time.Now()
	for _,ok := s.tokens[tok]; ok; tok = s.GetSessionToken() {}
	s.tokens[tok] = user

	session.Set("Login", tok)
	fmt.Println("Logged in!")

	ucpy := *user
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(ucpy)

	return 200, buf.String()
}

func (s *AEServer) HandleQuestionResponse(sess sessions.Session, params martini.Params, r *http.Request) (int, string) {
	id := bson.ObjectIdHex(params["id"])
	user := s.GetAuthedUser(sess)
	if user == nil {
		return 401, "{\"Message\":\"Not authorized to reply!\"}"
	}
	reply := new(Response)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(reply)
	if err != nil {
		return http.StatusBadRequest, "{\"Message\":\"Poorly formatted JSON\"}"
	}

	reply.ID = bson.NewObjectId()
	reply.Timestamp = time.Now()
	reply.Author = user.Username

	question,ok := s.questions.FindByID(id).(*Question)
	if !ok {
		return http.StatusForbidden, "{\"Message\":\"No such question!\"}"
	}
	question.Responses = append(question.Responses, reply)
	s.questions.Save(question)

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(reply)

	return 200, buf.String()
}

func (s *AEServer) HandleMe(session sessions.Session) (int, string) {
	return 200, "Nothing here"
}

func (s *AEServer) HandleVote(params martini.Params, session sessions.Session, r *http.Request) int {
	opt := params["opt"]
	if opt != "up" && opt != "down" {
		return 404
	}
	user := s.GetAuthedUser(session)
	if user == nil {
		return http.StatusUnauthorized
	}
	q := bson.ObjectIdHex(params["id"])
	question := s.questions.FindByID(q)
	if question == nil {
		return 404
	}

	return 200
}

func (s *AEServer) GetAuthedUser(sess sessions.Session) *User {
	//Verify user account or something
	login := sess.Get("Login")
	if login == nil {
		log.Printf("Not logged in!!")
		return nil
	}
	tok := login.(string)
	user, ok := s.tokens[tok]
	if !ok {
		log.Printf("Invalid cookie!")
		return nil
	}
	return user
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
