package main

import (
	"labix.org/v2/mgo/bson"
	"encoding/json"
	"bytes"
	"time"
	"io"
)

type Question struct {
	ID bson.ObjectId "_id,omitempty"
	Title string
	Author string
	Tags []string
	Upvotes []bson.ObjectId
	Downvotes []bson.ObjectId
	Timestamp time.Time
	LastEdit time.Time
	Body string
	Responses []*Response
	Comments []*Comment
}

func QuestionFromJson(r io.Reader) *Question {
	q := new(Question)
	dec := json.NewDecoder(r)
	err := dec.Decode(q)
	if err != nil {
		return nil
	}
	return q
}

func (q *Question) JsonBytes() []byte {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(q)
	return buf.Bytes()
}

func (q *Question) New() I {
	return new(Question)
}

func (q *Question) GetID() bson.ObjectId {
	return q.ID
}

func (q *Question) GetResponse(id bson.ObjectId) *Response {
	for _,v := range q.Responses {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func (q *Question) HasVoteBy(user bson.ObjectId) int {
	for _,v := range q.Upvotes {
		if v == user {
			return 1
		}
	}
	for _,v := range q.Downvotes {
		if v == user {
			return -1
		}
	}
	return 0
}

func (q *Question) Upvote(user bson.ObjectId) bool {
	switch q.HasVoteBy(user) {
	case 0:
		q.Upvotes = append(q.Upvotes, user)
		return true
	case 1:
		return false
	case -1:
		for i,v := range q.Downvotes {
			if v == user {
				q.Downvotes = append(q.Downvotes[:i], q.Downvotes[i+1:]...)
				q.Upvotes = append(q.Upvotes, user)
				return true
			}
		}
	}
	return false
}

func (q *Question) Downvote(user bson.ObjectId) bool {
	switch q.HasVoteBy(user) {
	case 0:
		q.Downvotes = append(q.Downvotes, user)
		return true
	case -1:
		return false
	case 1:
		for i,v := range q.Upvotes {
			if v == user {
				q.Upvotes = append(q.Upvotes[:i], q.Upvotes[i+1:]...)
				q.Downvotes = append(q.Downvotes, user)
				return true
			}
		}
	}
	return false
}

type Response struct {
	ID bson.ObjectId
	Author string
	Timestamp time.Time
	//Score Score
	Body string
	Comments []*Comment
}

func ResponseFromJson(r io.Reader) *Response {
	resp := new(Response)
	dec := json.NewDecoder(r)
	err := dec.Decode(resp)
	if err != nil {
		return nil
	}
	return resp
}

func (r *Response) AddComment(c *Comment) {
	r.Comments = append(r.Comments, c)
}

type Comment struct {
	ID bson.ObjectId
	Timestamp time.Time
	Author string
	Body string
	//Score Score
}

func CommentFromJson(r io.Reader) *Comment {
	comment := new(Comment)
	dec := json.NewDecoder(r)
	err := dec.Decode(comment)
	if err != nil {
		return nil
	}
	return comment
}

func (c *Comment) JsonBytes() []byte {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(c)
	return buf.Bytes()
}
