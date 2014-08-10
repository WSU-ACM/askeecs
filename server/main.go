package main

import (
	"github.com/whyrusleeping/askeecs/server/rest"
)

func main () {
	api := rest.Api{}
	api.ConnString = "localhost:27017"
	api.DBName = "askeecs"
	api.Run(":8080")
}
