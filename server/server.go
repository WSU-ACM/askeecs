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

func genRandString() string {
	buf := new(bytes.Buffer)
	io.CopyN(buf, rand.Reader, 32)
	return hex.EncodeToString(buf.Bytes())
}

func (s *AEServer) GetSessionToken() string {
	tok := genRandString()
	for _,ok := s.tokens[tok]; ok; tok = genRandString() {}
	return tok
}

func (s *AEServer) HandlePostQuestion(w http.ResponseWriter, r *http.Request, session sessions.Session) (int,string) {
	//Verify user account or something
	login := session.Get("Login")
	if login == nil {
		return 404, Message("Not Logged In!")
	}
	tok := login.(string)
	user, ok := s.tokens[tok]
	if !ok {
		return http.StatusBadRequest, Message("Invalid Cookie!")
	}

	q := QuestionFromJson(r.Body)
	if q == nil {
		return 404, Message("Poorly Formatted JSON.")
	}
	q.ID = bson.NewObjectId()
	q.Author = user.Username
	q.Timestamp = time.Now()

	err := s.questions.Save(q)
	if err != nil {
		log.Print(err)
		return http.StatusInternalServerError, Message("Failed to save question")
	}
	return 200, q.GetIdHex()
}

func (s *AEServer) HandleGetQuestions()(int,string) {
	q := s.questions.FindWhere(bson.M{})
	if q == nil {
		return 404,Message("Question not found.")
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

func (s *AEServer) HandleLogin(r *http.Request, params martini.Params, session sessions.Session) (int,string) {
	a := AuthFromJson(r.Body)
	if a == nil {
		time.Sleep(time.Second)
		return 404, Message("Login Failed")
	}

	users := s.users.FindWhere(bson.M{"username":a.Username})
	if len(users) == 0 {
		fmt.Println("User not found.")
		time.Sleep(time.Second)
		return 401, Message("Invalid Username or Password")
	}

	user, _ := users[0].(*User)
	if user.Password != a.Password {
		fmt.Println("Invalid password.")
		time.Sleep(time.Second)
		return http.StatusUnauthorized, Message("Invalid Username or Password.")
	}

	user.login = time.Now()
	tok := s.GetSessionToken()
	s.tokens[tok] = user
	session.Set("Login", tok)

	fmt.Println("Logged in!")
	return 200, string(user.JsonBytes())
}

func (s *AEServer) HandleQuestionComment(sess sessions.Session, params martini.Params, r *http.Request) (int, string) {
	id := bson.ObjectIdHex(params["id"])
	user := s.GetAuthedUser(sess)
	if user == nil {
		return 401, "{\"Message\":\"Not authorized to reply!\"}"
	}

	comment := CommentFromJson(r.Body)
	if comment == nil {
		return http.StatusBadRequest, Message("Poorly formatted JSON")
	}

	comment.Author = user.Username
	comment.Timestamp = time.Now()
	comment.ID = bson.NewObjectId()

	question,ok := s.questions.FindByID(id).(*Question)
	if !ok {
		return http.StatusForbidden, Message("No such question!")
	}
	question.Comments = append(question.Comments, comment)
	s.questions.Update(question)

	return 200, string(comment.JsonBytes())

}

func (s *AEServer) HandleQuestionResponse(sess sessions.Session, params martini.Params, r *http.Request) (int, string) {
	id := bson.ObjectIdHex(params["id"])
	user := s.GetAuthedUser(sess)
	if user == nil {
		return 401, "{\"Message\":\"Not authorized to reply!\"}"
	}

	reply := ResponseFromJson(r.Body)
	if reply == nil {
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
	s.questions.Update(question)

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(reply)

	return 200, buf.String()
}

func (s *AEServer) HandleResponseComment(sess sessions.Session, params martini.Params, r *http.Request) (int, string) {
	id := bson.ObjectIdHex(params["id"])
	user := s.GetAuthedUser(sess)
	if user == nil {
		return 401, "{\"Message\":\"Not authorized to reply!\"}"
	}

	comment := CommentFromJson(r.Body)
	if comment == nil {
		return http.StatusBadRequest, "{\"Message\":\"Poorly formatted JSON\"}"
	}
	comment.Author = user.Username
	comment.Timestamp = time.Now()
	comment.ID = bson.NewObjectId()

	question,ok := s.questions.FindByID(id).(*Question)
	if !ok {
		return http.StatusForbidden, "{\"Message\":\"No such question!\"}"
	}
	resp_id := params["resp"]
	resp := question.GetResponse(bson.ObjectId(resp_id))
	resp.AddComment(comment)

	s.questions.Update(question)

	return 200, string(comment.JsonBytes())
}

func (s *AEServer) HandleMe(session sessions.Session) (int, string) {
	return 200, "Nothing here"
}

func (s *AEServer) HandleVote(params martini.Params, session sessions.Session, r *http.Request) (int,string) {
	opt := params["opt"]
	if opt != "up" && opt != "down" {
		return http.StatusMethodNotAllowed,"{\"Message\":\"Invalid vote type\"}"
	}
	user := s.GetAuthedUser(session)
	if user == nil {
		return http.StatusUnauthorized, "{\"Message\":\"Not logged in!\"}"

	}
	q := bson.ObjectIdHex(params["id"])
	question,ok := s.questions.FindByID(q).(*Question)
	if question == nil || !ok {
		return 404, "{\"Message\":\"No such question!\"}"
	}
	switch opt {
		case "up":
			if question.Upvote(user._id) {
				s.questions.Update(question)
			}
		case "down":
			if question.Downvote(user._id) {
				s.questions.Update(question)
			}
	}

	return 200, string(question.JsonBytes())
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

func Message(s string) string {
	return fmt.Sprintf("{\"Message\":\"%s\"}", s)
}
