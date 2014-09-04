package rest

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"errors"
	. "github.com/visionmedia/go-debug"
)

var db_debug = Debug("askeecs:db")

var ErrorNotFound = errors.New("No documents found!")
var ErrorNullResponse = errors.New("Got back null response from mgo.")

type Database struct {
	db *mgo.Database
	collections map[string]*Collection
}

func NewDatabase(host string, db_name string) *Database {
	s,err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	mdb := s.DB(db_name)

	dbs := new(Database)
	dbs.db = mdb
	dbs.collections = make(map[string]*Collection)
	return dbs
}

func (db *Database) Collection(name string, typ I) *Collection {
	c,ok := db.collections[name]
	if ok {
		return c
	}

	c = new(Collection)
	c.col = db.db.C(name)
	c.template = typ
	db.collections[name] = c
	return c
}

type I interface {
	GetID() bson.ObjectId
	New() I
}

type Collection struct {
	col *mgo.Collection
	cache map[string]I
	template I
}

func (c *Collection) Save(doc I) error {
	//TODO: handle errors?
	db_debug("Saving document.")
	err := c.col.Insert(doc)
	return err
}

func (c *Collection) Update(doc I) error {
	db_debug("Updating Document.")
	err := c.col.UpdateId(doc.GetID(), doc)
	return err
}

func (c *Collection) FindByID(ID bson.ObjectId) I {
	q := c.col.FindId(ID)

	if q == nil {
		db_debug("%s", ErrorNullResponse)
		return nil
	}

	cnt,err := q.Count()

	if err != nil {
		db_debug("%s", err)
		return nil
	}
	if cnt < 1 {
		db_debug("%s", ErrorNotFound)
		return nil
	}
	obj := c.template.New()
	q.One(obj)
	return obj
}

func (c *Collection) FindWhere(match bson.M) []I {
	db_debug("%s", match)
	q := c.col.Find(match)
	if q == nil {
		db_debug("%s", ErrorNullResponse)
		return nil
	}

	n,err := q.Count()
	if err != nil {
		db_debug("%s", err)
		return nil
	}
	var out []I

	if n == 0 {
		db_debug("Nothing matched the query...")
		return out
	}

	i := q.Iter()
	v := c.template.New()
	for i.Next(v) {
		out = append(out,v)
		v = c.template.New()
	}
	return out
}

func (c *Collection) FindOneWhere(match bson.M) I {
	q := c.col.Find(match)
	if q == nil {
		db_debug("%s", ErrorNullResponse)
		return nil
	}
	cnt,err := q.Count()
	if err != nil {
		db_debug("%s", err)
		return nil
	}
	if cnt < 1 {
		db_debug("%s", ErrorNotFound)
		return nil
	}
	obj := c.template.New()
	q.One(obj)
	return obj
}

