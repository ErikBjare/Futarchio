package db

import (
	"appengine"
	"appengine/datastore"
	"time"
)

// UserCreated - Serves as a base entity for things created by users such as comments, polls and predictions.
type UserCreated struct {
	Creator *datastore.Key `json:"creator"`
	Created time.Time      `json:"created"`
}

func newUserCreated(creator *datastore.Key) UserCreated {
	return UserCreated{
		Creator: creator,
		Created: time.Now(),
	}
}

// Post - Serves as a base entity for things created by users with a title and description.
type Post struct {
	UserCreated
	Title       string `json:"title"`
	Description string `json:"description"`
}

// NewPost - Creates a new Post
func NewPost(title, desc string, creator *datastore.Key) Post {
	return Post{
		UserCreated: newUserCreated(creator),
		Title:       title,
		Description: desc,
	}
}

// Comment is simply a basic comment, can be voted and commented on.
// TODO: Far from done
// How to make fetching nested comments efficient?
// Option 1: Keep comment-able entity as parent, have parent as
// Should probably be the ancestor of a parent comment/post.
type Comment struct {
	UserCreated
	Thread *datastore.Key
	Text   string
}

// Entity
//
// Maybe a decent candidate for storing entities together with their keys
type Entity struct {
	key *datastore.Key
}

func NewEntity(c *appengine.Context, kind string, parent *datastore.Key) *Entity {
	return &Entity{key: datastore.NewIncompleteKey(*c, kind, parent)}
}

func (e *Entity) Key() (*datastore.Key, error) {
	return e.key, nil
}

func (e *Entity) SetKey(key *datastore.Key) {
	e.key = key
}

// Votable is a base entity for things that can be voted on
// TODO: Store votes, preventing double-voting
// Should just use a simpler votereceipt.
// Can 'votable' handle race-conditions?
// Perhaps a better solution would be to simply do:
//	q := datastore.NewQuery("PostVote").Filter("Key =", key)
//	up := q.Filter("IsUp =", true).Count()
//	down := q.Filter("IsUp =", false).Count()
type Votable struct {
	Entity
}

type Votes struct {
	Up   int
	Down int
}

func (v *Votable) CountVotes(c appengine.Context) (*Votes, error) {
	key, err := v.Key()
	if err != nil {
		return nil, err
	}
	q := datastore.NewQuery("PostVote").Filter("Key =", key)
	up, err := q.Filter("IsUp =", true).Count(c)
	if err != nil {
		return nil, err
	}
	down, err := q.Filter("IsUp =", false).Count(c)
	if err != nil {
		return nil, err
	}

	return &Votes{Up: up, Down: down}, nil
}
