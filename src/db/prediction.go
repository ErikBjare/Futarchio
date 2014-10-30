package db

import (
	"appengine"
	"appengine/datastore"
	"time"
)

// Represents a statement which people can make predictions o, create with poll initializers
// TODO: Finalization condition (time/expiry, vote, etc.)
type Claim struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Creator     string    `json:"creator"`
	Created     time.Time `json:"created"`
	// Type can be one of ["CredenceClaim"]
	Type string `json:"type"`
}

// Creates a new poll.
//
// Should rarely be used, use specialized poll constructors instead.
func newClaim(title, desc string, creator string) Claim {
	return Claim{
		Title:       title,
		Description: desc,
		Creator:     creator,
		Created:     time.Now(),
	}
}

// Returns the current standings of the claim
// TODO: Calculate average, mean and stddev
func (p *Claim) Stats(c appengine.Context, claimkey *datastore.Key) map[string]float32 {
	c.Errorf("Unimplemented: Standings")
	var votes []Vote
	q := datastore.NewQuery("Prediction").Filter("Claim =", claimkey)
	_, err := q.GetAll(c, &votes)
	if err != nil {
		// TODO: Better error handling
		panic(err)
	}
	// TODO: Return actual value
	return map[string]float32{}
}

// Creates a claim based on credence
func NewCredenceClaim(title, desc string, creator string) Claim {
	p := newClaim(title, desc, creator)
	p.Type = "CredenceClaim"
	return p
}

// A Prediction on a claim
type Prediction struct {
	Claim *datastore.Key `json:"claimid"`

	// The username of the predictor
	User string `json:"userid"`

	// The credence assigned by the predictor, a value in the open interval (0-1)
	Credence float32 `json:"credence"`

	// The date and time of creation
	Created time.Time `json:"created"`
}
