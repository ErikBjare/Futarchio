package db

import (
	//	"fmt"
	"appengine/datastore"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"math/rand"
	"strconv"
	"time"
)

/*
   User-related models
*/

type User struct {
	Username string    `json:"username"`
	Password []byte    `json:"-"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
}

type Auth struct {
	// Should always be an ancestor of User
	Key     string    `json:"key"`
	Expires time.Time `json:"expires"`
}

func NewUser(username string, password string, name string, email string) *User {
	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		log.Fatal(err)
	}

	return &User{
		Username: username,
		Password: hashed_pass,
		Name:     name,
		Email:    email,
		Created:  time.Now(),
	}
}

func NewAuth() Auth {
	auth_bytes := sha256.Sum256([]byte(strconv.Itoa(rand.Int())))
	authkey := base64.StdEncoding.EncodeToString([]byte(auth_bytes[:]))
	return Auth{
		Key:     authkey,
		Expires: time.Now().Add(24 * time.Hour),
	}
}

func (u User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		log.Println(err)
	}
	return !(err != nil)
}

/*
   Poll-related models
*/

type Poll struct {
	// Represents a base poll, needs to be filled by poll-type initializers
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Creator     *datastore.Key `json:"creator"`
	// Type can be one of "YesNoPoll", "CredencePoll", "MultipleChoicePoll", "AllocationPoll"
	Type    string   `json:"type"`
	Choices []string `json:"choices"`
}

func NewPoll(title, desc string, creator *datastore.Key) Poll {
	return Poll{
		Title:       title,
		Description: desc,
		Creator:     creator,
	}
}

func (p Poll) AddChoice(name string) {
	p.Choices = append(p.Choices, name)
}

func NewYesNoPoll(title, desc string, creator *datastore.Key) Poll {
	p := NewPoll(title, desc, creator)
	p.Type = "YesNoPoll"
	p.Choices = []string{"yes", "no"}
	return p
}

func MultichoicePoll(title, desc string, creator *datastore.Key, choices []string) Poll {
	p := NewPoll(title, desc, creator)
	p.Type = "MultichoicePoll"
	p.Choices = choices
	return p
}

type Vote struct {
	// Should always have a Poll as parent
	Weights map[string]float32 `json:"weights" datastore:"weights"`
	Key     string             // Optional, never both Creator and Key
	Creator *datastore.Key     // Optional, never both Creator and Key
}

func NewVote(choice map[string]float32) Vote {
	return Vote{
		Weights: choice,
	}
}

func NewYesNoVote(yes bool) Vote {
	if yes {
		return NewVote(map[string]float32{"yes": 1})
	} else {
		return NewVote(map[string]float32{"no": 1})
	}
}

func SumVotes(vs []Vote) map[string]float32 {
	weights := map[string]float32{"yes": 0, "no": 0}
	for _, v := range vs {
		for k, f := range v.Weights {
			weights[k] += f
		}
	}
	return weights
}
