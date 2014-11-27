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
// Option 1: Keep comment-able entity as parent, keep parent _post_ key in a variable.
// Should probably be the ancestor of a parent comment/post.
type Comment struct {
	UserCreated
	Thread *datastore.Key
	Text   string
}

// Votable is a base entity for things that can be voted up or down reddit/HN/SE-style
// TODO: Store votes, preventing double-voting
type Votable struct {
}

type Votes struct {
	Up   int
	Down int
}

func CountVotes(c appengine.Context, key *datastore.Key) (*Votes, error) {
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
