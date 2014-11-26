package db

import (
	"appengine"
	"appengine/datastore"
)

// Statement - Represents a statement which people can make predictions on, create with poll initializers
type Statement struct {
	Post
	// TODO: Finalization condition/type (time/expiry, vote, etc.)
}

// Creates a new statement.
//
// Should rarely be used, use specialized poll constructors instead.
func NewStatement(title, desc string, creator *datastore.Key) Statement {
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

// A Prediction on a statement
type Prediction struct {
	UserCreated

	Statement *datastore.Key `json:"statement"`
	// The credence assigned by the predictor, a value in the open interval (0-1)
	Credence float32 `json:"credence"`
}

func NewPrediction(userkey *datastore.Key, stmtkey *datastore.Key, credence float32) Prediction {
	return Prediction{
		UserCreated: newUserCreated(userkey),
		Statement:   stmtkey,
		Credence:    credence,
	}
}
