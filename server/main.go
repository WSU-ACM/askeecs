package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/sessions"
	//"github.com/shykes/spdy-go"

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

	m.Get("/q/:id", s.HandleGetQuestion)
	m.Get("/q/:id/vote/:opt", s.HandleVote)
	m.Post("/q", s.HandlePostQuestion)
	m.Get("/q", s.HandleGetQuestions)
	m.Post("/login", s.HandleLogin)
	m.Post("/register", s.HandleRegister)
	m.Post("/logout", s.HandleLogout)
	m.Post("/me", s.HandleMe);
	m.Run()
	//spdy.ListenAndServeTCP(":3000", m)
}
