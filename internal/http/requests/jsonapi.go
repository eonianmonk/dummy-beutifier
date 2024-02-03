package requests

import "time"

type Article struct {
	Id        int        `jsonapi:"primary,post"`
	Title     string     `jsonapi:"attr,title"`
	Body      string     `jsonapi:"attr,body"`
	CreatedAt time.Time  `jsonapi:"attr, created_at"`
	Author    string     `jsonapi:"attr,author"`
	Comments  []*Comment `jsonapi:"relation,comments"`
}

type Comment struct {
	Id        int       `jsonapi:"primary,comment"`
	Author    string    `jsonapi:"attr,author"`
	Body      string    `jsonapi:"attr,body"`
	CreatedAt time.Time `jsonapi:"attr, created_at"`
}
