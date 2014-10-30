package db

import (
	"appengine"
	"appengine/datastore"
	"time"
)

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
