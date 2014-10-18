package db

import (
	"appengine/datastore"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// A user
type User struct {
	Username string    `json:"username"`
	Password []byte    `json:"-"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
}

// Creates a new user
func NewUser(username string, password string, name string, email string) User {
	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		log.Fatal(err)
	}

	return User{
		Username: username,
		Password: hashed_pass,
		Name:     name,
		Email:    email,
		Created:  time.Now(),
	}
}

// Checks the password for a given user
func (u User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		log.Println(err)
	}
	return !(err != nil)
}

// Represents an authentication key
//
// Should always be an ancestor of User in the datastore
type Auth struct {
	// Should always be an ancestor of User
	Key     string    `json:"key"`
	Expires time.Time `json:"expires"`
}

func NewAuth() Auth {
	auth_bytes := sha256.Sum256([]byte(strconv.Itoa(rand.Int())))
	authkey := base64.StdEncoding.EncodeToString([]byte(auth_bytes[:]))
	return Auth{
		Key:     authkey,
		Expires: time.Now().Add(24 * time.Hour),
	}
}

// Represents a base poll, create with poll initializers
type Poll struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Creator     string    `json:"creator"`
	Created     time.Time `json:"created"`
	// Type can be one of "YesNoPoll", "CredencePoll", "MultipleChoicePoll", "AllocationPoll"
	Type    string   `json:"type"`
	Choices []string `json:"choices"`
}

// Creates a new poll.
//
// Should rarely be used, use specialized poll constructors instead.
func newPoll(title, desc string, creator string) Poll {
	return Poll{
		Title:       title,
		Description: desc,
		Creator:     creator,
		Created:     time.Now(),
	}
}

// Adds a choice to a poll
func (p Poll) AddChoice(name string) {
	p.Choices = append(p.Choices, name)
}

// Creates a yes/no poll
func NewYesNoPoll(title, desc string, creator string) Poll {
	p := newPoll(title, desc, creator)
	p.Type = "YesNoPoll"
	p.Choices = []string{"yes", "no"}
	return p
}

// Creates a multiple choice poll
func MultichoicePoll(title, desc string, creator string, choices []string) Poll {
	p := newPoll(title, desc, creator)
	p.Type = "MultichoicePoll"
	p.Choices = choices
	return p
}

// A vote
//
// Should always have a Poll as parent
type Vote struct {
	// The weights of the vote
	//
	// Keys are options, values are how much of the vote is put on each choice
	Weights map[string]float32 `json:"weights" datastore:"weights"`

	// The username of the voter
	//
	// Optional, never both Creator and Key
	Creator string `json:"creator"`

	// Represents a public key when a vote is made anonymously
	//
	// Optional, never both Creator and Key
	Key string `json:"key"`
}

// A receipt that the user has voted
//
// Only useful if voter voted anonymously
type VoterReceipt struct {
	Poll *datastore.Key
	User *datastore.Key
}

// Creates a new vote
//
// Is user == nil, then vote anonymously and return the private key
func newVote(choice map[string]float32, user User, anon bool) (Vote, string) {
	var private_key string
	vote := Vote{
		Weights: choice,
	}
	if anon {
		vote.Creator = user.Username
	} else {
		private_key = user.Username + "#" + strconv.Itoa(rand.Int())
		hash := sha256.Sum256([]byte(private_key))
		vote.Key = base64.StdEncoding.EncodeToString([]byte(hash[:]))
	}
	return vote, private_key
}

// Creates new Vote for a Yes or No poll.
//
// If user == nil, then vote anonymously
func NewYesNoVote(yes bool, user User, anon bool) (Vote, string) {
	var choice map[string]float32
	if yes {
		choice = map[string]float32{"yes": 1}
	} else {
		choice = map[string]float32{"no": 1}
	}
	return newVote(choice, user, anon)
}

// Sums a collection of votes
//
// TODO: Vote normalization
func SumVotes(vs []Vote) map[string]float32 {
	weights := map[string]float32{"yes": 0, "no": 0}
	for _, v := range vs {
		for k, f := range v.Weights {
			weights[k] += f
		}
	}
	return weights
}
