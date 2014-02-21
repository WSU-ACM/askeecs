# Ask EECS

A stackoverflow like site for the WSU EECS Community.

![Ask EECS](http://i.imgur.com/MvlS2yt.png)
*I'm sorry for the spelling errors in the screenshot, I edited it in the developer console*

## How To Use

 1. Install [Go](http://golang.org)
 2. Make sure your `GOPATH` is configured properly.
 3. Install [Martini](http://martini.codegangsta.io):
  - `go get github.com/codegangsta/martini`
  - `go get github.com/martini-contrib/sessions`
 4. Install [MongoDB](http://www.mongodb.org).
 5. Install the [Mongo bindings for Go](http://labix.org/mgo):
  - `go get labix.org/v2/mgo`
  - `go get labix.org/v2/mgo/bson`
 6. Run `make`!
 7. Spin up MongoDB:
  - `mkdir data && mongodb --dbpath data`
 8. Finally, run `./askeecs`.

By default, martini serves to localhost at port 3000.

## Contributors:
- Jeromy Johnson
- Travis Person
