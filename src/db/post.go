package db

import (
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

// Comment - A basic comment
// TODO: Far from done
// Should probably be the ancestor of a parent comment/post.
type Comment struct {
	UserCreated
	Thread string
	Text   string
}
