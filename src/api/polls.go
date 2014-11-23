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

type VoteResponse struct {
	db.Vote
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
		Filter(authFilter).
		Doc("create a poll").
		Operation("createPoll").
		Reads(db.Poll{}))
	ws.Route(ws.POST("/{pollid}/vote").To(p.vote).
		Doc("vote on a poll").
		Operation("vote").
		Filter(authFilter).
		Param(ws.PathParameter("pollid", "Id of poll to vote on").DataType("string")).
		Reads(VoteRequest{}).
		Writes(db.VoteReceipt{}))

	// The following two routes use the same method
	ws.Route(ws.GET("/votes").To(p.getVotes).
		Doc("get the latest votes for all polls").
		Operation("getVotes").
		Writes([]VoteResponse{}))
	ws.Route(ws.GET("/{pollid}/votes").To(p.getVotes).
		Doc("get the votes for given pollid").
		Operation("getPollVotes").
		Param(ws.PathParameter("pollid", "id of poll to get votes for").DataType("string")).
		Writes([]VoteResponse{}))

	// TODO: Move these two to User endpoint?
	ws.Route(ws.GET("/myvotereceipts").To(p.getMyVotereceipts).
		Filter(authFilter).
		Doc("get the current users votereceipts").
		Operation("getMyVotereceipts").
		Writes([]db.VoteReceipt{}))

	restful.Add(ws)
}

func fetchVotes(c appengine.Context, user *datastore.Key) (*[]db.Vote, error) {
	// TODO: Fetch all of the users votereceipts, use the keys to find private votes

	// Fetch all votes which are public
	var votes []db.Vote
	q := datastore.NewQuery("Vote")
	if user != nil {
		q = q.Filter("Creator =", user)
	}
	_, err := q.GetAll(c, &votes)
	if err != nil {
		return nil, err
	}

	return &votes, nil
}

func (p PollApi) getMyVotereceipts(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	_, userkey := auth(c, r)

	var votereceipts []db.VoteReceipt
	q := datastore.NewQuery("VoteReceipt").Filter("User =", userkey)
	_, err := q.GetAll(c, &votereceipts)
	if err != nil {
		c.Errorf(err.Error())
		respondError(w, 500, "")
		return
	}

	respondMany(w, votereceipts)
}

func (p PollApi) getVotes(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	pollid := r.PathParameter("pollid")
	c.Infof(pollid)

	var votes []VoteResponse
	q := datastore.NewQuery("Vote")
	if pollid != "" {
		pollkey, err := datastore.DecodeKey(pollid)
		if err != nil {
			c.Errorf("%v", err)
		} else {
			q = q.Filter("Poll =", pollkey)
		}
	}
	_, err := q.GetAll(c, &votes)
	if err != nil {
		c.Errorf(err.Error())
		respondError(w, 500, err.Error())
		return
	}

	for i := range votes {
		votes[i].Weights = votes[i].Vote.Weights()
	}

	c.Infof("%v", votes)
	respondMany(w, votes)
}

func (p PollApi) vote(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	_, userkey := auth(c, r)

	pollkey, err := datastore.DecodeKey(r.PathParameter("pollid"))
	if err != nil {
		c.Errorf("%v", err)
	}

	var poll db.Poll
	err = datastore.Get(c, pollkey, &poll)
	if err != nil {
		c.Errorf("%v", err)
	}

	// Check if user has already voted
	// TODO: Speed this process up by memcaching the user-poll relationship
	var votereceipts []db.VoteReceipt
	q := datastore.NewQuery("VoteReceipt").Filter("Poll =", pollkey).Filter("User =", userkey)
	vrkeys, err := q.GetAll(c, &votereceipts)
	if err != nil {
		c.Errorf("%v", err)
	}
	if len(vrkeys) > 0 {
		errstr := "attempted to double vote"
		c.Infof(errstr)
		respondError(w, 500, errstr)
		return
	}

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
		vote, votereceipt, private_key = db.NewYesNoVote(pollkey, userkey, yes, db.Public)
	} else {
		c.Errorf("%s", poll)
		c.Errorf("invalid polltype: %s", poll.Type)
		respondError(w, 500, "invalid polltype: "+poll.Type)
		return
	}

	// Save vote and votereceipt
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
	// TODO: Pagination
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

	_, userkey := auth(c, r)

	var poll db.Poll
	err := r.ReadEntity(&poll)
	if err != nil {
		c.Errorf("Error")
	}

	c.Debugf("Created poll: %s", poll)
	if poll.Type == "YesNoPoll" {
		poll := db.NewYesNoPoll(poll.Title, poll.Description, userkey)
		datastore.Put(c, datastore.NewIncompleteKey(c, "Poll", nil), &poll)
		c.Infof("Saved poll!")
	}
}
