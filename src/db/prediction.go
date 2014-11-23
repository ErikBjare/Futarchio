package db

import (
	"appengine"
	"appengine/datastore"
)

// Statement - Represents a statement which people can make predictions o, create with poll initializers
type Statement struct {
	Post
	// Type can be one of ["CredenceClaim"]
	Type string `json:"type"`
	// TODO: Finalization condition/type (time/expiry, vote, etc.)
}

// Creates a new statement.
//
// Should rarely be used, use specialized poll constructors instead.
func newStatement(title, desc string, creator *datastore.Key) Statement {
	return Statement{
		Post: NewPost(title, desc, creator),
	}
}

// Stats - Returns the current standings of the claim
// TODO: Calculate average, mean and stddev
func (p *Statement) Stats(c appengine.Context, claimkey *datastore.Key) map[string]float32 {
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

// NewCredenceStatement - Creates a statement where predictions are assigned a credence score by their predictors. The standard type of statement.
// TODO: Rename to something more intuitive
func NewCredenceStatement(title, desc string, creator *datastore.Key) Statement {
	p := newStatement(title, desc, creator)
	p.Type = "CredenceClaim"
	return p
}

// A Prediction on a statement
type Prediction struct {
	Post

	Claim *datastore.Key `json:"claimid"`
	// The credence assigned by the predictor, a value in the open interval (0-1)
	Credence float32 `json:"credence"`
}
