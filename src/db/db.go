package db

import (
	"appengine"
	"appengine/datastore"
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"log"
	"math/rand"
	"strconv"
	"strings"
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
		Username: strings.ToLower(username),
		Password: hashed_pass,
		Name:     name,
		Email:    strings.ToLower(email),
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
func (p *Poll) AddChoice(name string) {
	p.Choices = append(p.Choices, name)
}

// Returns the current standings of the poll
func (p *Poll) Weights(c appengine.Context, pollkey *datastore.Key) map[string]float32 {
	var votes []Vote
	q := datastore.NewQuery("Vote").Filter("Poll =", pollkey)
	_, err := q.GetAll(c, &votes)
	if err != nil {
		// TODO: Better error handling
		panic(err)
	}
	return SumVotes(votes)
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
type Vote struct {
	Poll *datastore.Key `json:"pollid"`

	// The weights of the vote, a map[string]float32
	//
	// Keys are options, values are how much of the vote is put on each choice
	EncodedWeights []byte `json:"-"`

	// The username of the voter
	//
	// Optional, never both Creator and Key
	Creator *datastore.Key `json:"creatorid"`

	// Represents a public key when a vote is made anonymously
	//
	// Optional, never both Creator and Key
	Key string `json:"key"`

	// The time and date of creation
	Created time.Time `json:"created"`
}

func (v *Vote) Weights() map[string]float32 {
	reader := bytes.NewReader(v.EncodedWeights)
	dec := gob.NewDecoder(reader)
	var weights map[string]float32
	err := dec.Decode(&weights)
	if err != nil {
		panic(err)
	}
	return weights
}

func (v *Vote) SetWeights(w map[string]float32) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(w)
	if err != nil {
		return err
	}
	v.EncodedWeights = buffer.Bytes()
	return nil
}

// A receipt that the user has voted
//
// Only useful if voter voted anonymously
type VoteReceipt struct {
	Poll *datastore.Key `json:"pollid"`
	User *datastore.Key `json:"userid"`
	// If vote is private, store key here
	Key string `json:"key"`
}

const (
	Public    = 0
	Private   = 5
	Anonymous = 10
)

// Creates a new vote
//
// Is user == nil, then vote anonymously and return the private key
// Privacy is either Public (0), Private (5) or Anonymous (10)
// TODO: Add entropy to private key, use bcrypt?
func newVote(pollkey *datastore.Key, userkey *datastore.Key, choice map[string]float32, privacy int) (Vote, VoteReceipt, string) {
	private_key := userkey.Encode() + "#" + strconv.Itoa(rand.Int())

	hash := sha256.Sum256([]byte(private_key))
	vote := Vote{
		Poll:    pollkey,
		Key:     base64.StdEncoding.EncodeToString([]byte(hash[:])),
		Created: time.Now(),
	}
	vote.SetWeights(choice)
	if len(vote.EncodedWeights) == 0 {
		log.Println("len of encoded weights was 0")
	}

	votereceipt := VoteReceipt{
		Poll: pollkey,
		User: userkey,
	}

	// Store user in vote if public
	if privacy == Public {
		vote.Creator = userkey
	}

	// If vote is private or public, store private key in receipt
	if privacy <= Private {
		votereceipt.Key = private_key
	}

	return vote, votereceipt, private_key
}

// Creates new Vote for a Yes or No poll.
//
// If user == nil, then vote anonymously
func NewYesNoVote(pollkey *datastore.Key, userkey *datastore.Key, yes bool, privacy int) (Vote, VoteReceipt, string) {
	var choice map[string]float32
	if yes {
		choice = map[string]float32{"yes": 1}
	} else {
		choice = map[string]float32{"no": 1}
	}
	return newVote(pollkey, userkey, choice, privacy)
}

// Sums a collection of votes
//
// TODO: Vote normalization
// TODO: Make it work not exclusively for yes and no
func SumVotes(vs []Vote) map[string]float32 {
	weights := map[string]float32{"yes": 0, "no": 0}

	// Stores the total size of the vote weights for normalization
	var votesum float32 = 0.0
	for _, v := range vs {
		w := v.Weights()
		for _, f := range w {
			votesum += f
		}
		for k, f := range w {
			weights[k] += f / votesum
		}
		votesum = 0.0
	}
	return weights
}
