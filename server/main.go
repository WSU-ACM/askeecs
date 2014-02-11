package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/sessions"

	"io/ioutil"
)

func main() {
	s := NewServer()
	m := martini.Classic()


	secret,err := ioutil.ReadFile(".secret")
	if err != nil {
		panic(err)
	}
	store := sessions.NewCookieStore(secret)
	m.Use(sessions.Sessions("ask_eecs_auth_session", store))

	m.Get("/q", s.HandleGetQuestions)
	m.Post("/q", s.HandlePostQuestion)
	m.Get("/q/:id", s.HandleGetQuestion)
	m.Put("/q/:id", s.HandleEditQuestion)
	m.Get("/q/:id/vote/:opt", s.HandleVote)
	m.Post("/q/:id/response", s.HandleQuestionResponse)
	m.Post("/q/:id/response/:resp/comment", s.HandleResponseComment)
	m.Post("/q/:id/comment", s.HandleQuestionComment)

	m.Post("/login", s.HandleLogin)
	m.Post("/register", s.HandleRegister)
	m.Post("/logout", s.HandleLogout)
	m.Post("/me", s.HandleMe);
	m.Run()
}
