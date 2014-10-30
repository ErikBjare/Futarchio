package db

import (
	"appengine/datastore"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// A vote
type Vote struct {
	Poll *datastore.Key `json:"pollid"`

	// The weights of the vote, a map[string]float32
	//
	// Keys are options, values are how much of the vote is put on each choice
	encodedWeights []byte `json:"-"`

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
	reader := bytes.NewReader(v.encodedWeights)
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
	v.encodedWeights = buffer.Bytes()
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
	if len(vote.encodedWeights) == 0 {
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
