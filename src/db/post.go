package db

import (
	"time"
)

// UserCreated - Serves as a base entity for things created by users such as comments, polls and predictions.
type UserCreated struct {
	Creator string    `json:"creator"`
	Created time.Time `json:"created"`
}

func newUserCreated(creator string) UserCreated {
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
func NewPost(title, desc, creator string) Post {
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
