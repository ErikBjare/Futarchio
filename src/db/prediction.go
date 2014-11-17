package db

import (
	"appengine"
	"appengine/datastore"
)

// Claim - Represents a statement which people can make predictions o, create with poll initializers
// TODO: Rename to statement?
type Claim struct {
	Post
	// Type can be one of ["CredenceClaim"]
	Type string `json:"type"`
	// TODO: Finalization condition/type (time/expiry, vote, etc.)
}

// Creates a new poll.
//
// Should rarely be used, use specialized poll constructors instead.
func newClaim(title, desc string, creator string) Claim {
	return Claim{
		Post: NewPost(title, desc, creator),
	}
}

// Stats - Returns the current standings of the claim
// TODO: Calculate average, mean and stddev
func (p *Claim) Stats(c appengine.Context, claimkey *datastore.Key) map[string]float32 {
	c.Errorf("Unimplemented: Claim.Stats")
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

// NewCredenceClaim - Creates a claim based on credence
// TODO: Rename to something more intuitive
func NewCredenceClaim(title, desc string, creator string) Claim {
	p := newClaim(title, desc, creator)
	p.Type = "CredenceClaim"
	return p
}

// A Prediction on a claim
type Prediction struct {
	Post

	Claim *datastore.Key `json:"claimid"`
	// The credence assigned by the predictor, a value in the open interval (0-1)
	Credence float32 `json:"credence"`
}
