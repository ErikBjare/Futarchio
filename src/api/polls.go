package api

import (
	//	"fmt"
	"appengine"
	"appengine/datastore"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
)

type PollApi Api

type PollResponse struct {
	// Id should be an encoded datastore.Key
	Id string `json:"id"`
	db.Poll
	Weights map[string]float32 `json:"weights"`
}

type VoteRequest struct {
	Weights map[string]float32 `json:"weights"`
}

func (p PollApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path("/api/0/polls").
		Doc("Polls").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(p.getLatest).
		Doc("get the latest polls").
		Operation("getLatest").
		Writes([]PollResponse{}))
	ws.Route(ws.POST("").To(p.createPoll).
		Filter(basicAuthenticate).
		Doc("create a poll").
		Operation("createPoll").
		Reads(db.Poll{}))
	ws.Route(ws.POST("/{pollid}/vote").To(p.vote).
		Doc("vote on a poll").
		Operation("vote").
		Filter(basicAuthenticate).
		Param(ws.PathParameter("pollid", "Id of poll to vote on").DataType("string")).
		Reads(VoteRequest{}).
		Writes(db.VoteReceipt{}))

	restful.Add(ws)
}

func (p PollApi) vote(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	user, _ := auth(c, r)

	pollkey, err := datastore.DecodeKey(r.PathParameter("pollid"))
	if err != nil {
		c.Errorf("%v", err)
	}

	var poll db.Poll
	err = datastore.Get(c, pollkey, &poll)
	if err != nil {
		c.Errorf("%v", err)
	}

	// TODO: Check if user has already voted, speed this process up by memcaching the user-poll relationship

	// Get incoming JSON data specifying poll
	var vote_req VoteRequest
	err = r.ReadEntity(&vote_req)
	if err != nil {
		c.Errorf("%v", err)
	}

	// Create vote & vote receipt
	var vote db.Vote
	var votereceipt db.VoteReceipt
	var private_key string
	if poll.Type == "YesNoPoll" {
		var yes bool
		if vote_req.Weights["yes"] == 1 {
			yes = true
		} else {
			yes = false
		}
		vote, votereceipt, private_key = db.NewYesNoVote(pollkey, *user, yes, db.Public)
	} else {
		c.Errorf("%s", poll)
		c.Errorf("invalid polltype: %s", poll.Type)
		respondError(w, 500, "invalid polltype: "+poll.Type)
		return
	}

	// TODO: Save vote and votereceipt
	votekey := datastore.NewIncompleteKey(c, "Vote", nil)
	votekey, err = datastore.Put(c, votekey, &vote)
	if err != nil {
		c.Errorf("%v", err)
	}

	votereceiptkey := datastore.NewIncompleteKey(c, "VoteReceipt", nil)
	votereceiptkey, err = datastore.Put(c, votereceiptkey, &votereceipt)
	if err != nil {
		c.Errorf("%v", err)
	}

	// Return votereceipt, with private key
	votereceipt.Key = private_key
	respondOne(w, votereceipt)
}

func (p PollApi) getLatest(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	q := datastore.NewQuery("Poll").Limit(20).Order("-Created")

	var polls []db.Poll
	keys, err := q.GetAll(c, &polls)
	if err != nil {
		c.Errorf(err.Error())
	}

	if len(keys) == 0 {
		respondMany(w, []db.Poll{})
		return
	}

	polls_response := make([]PollResponse, len(polls))
	for i, _ := range polls {
		polls_response[i].Poll = polls[i]
		polls_response[i].Id = keys[i].Encode()
		polls_response[i].Weights = polls[i].Weights(c, keys[i])
	}

	c.Debugf("%s", polls_response)
	respondMany(w, polls_response)
}

func (p PollApi) createPoll(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	user, _ := auth(c, r)

	var poll db.Poll
	err := r.ReadEntity(&poll)
	if err != nil {
		c.Errorf("Error")
	}

	c.Debugf("Created poll: %s", poll)
	if poll.Type == "YesNoPoll" {
		poll := db.NewYesNoPoll(poll.Title, poll.Description, user.Username)
		datastore.Put(c, datastore.NewIncompleteKey(c, "Poll", nil), &poll)
		c.Infof("Saved poll!")
	}
}
