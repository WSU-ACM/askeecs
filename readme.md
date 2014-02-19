#Ask EECS

A stackoverflow like site for the WSU EECS Community.

![Ask EECS](http://i.imgur.com/MvlS2yt.png)
*I'm sorry for the spelling errors in the screenshot, I edited it in the developer console*

##How To Use

- Install Go [http://golang.org](Go)
- Make sure your gopath is configured properly
- Install Martini:
  - `go get github.com/codegangsta/martini`
  - `go get github.com/martini-contrib/sessions`
- Install MongoDB (using your package manager)
- Install the mongo bindings for go:
  - `go get labix.org/v2/mgo`
  - `go get labix.org/v2/mgo/bson`
- run make!
- Spin up mongo db:
  - `mkdir data && mongodb --dbpath data`
- Finally, run `./askeecs`
- By default, martini serves to localhost at port 3000

##Contributors:
Jeromy Johnson
Travis Person
