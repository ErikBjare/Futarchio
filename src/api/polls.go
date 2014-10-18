package api

import (
	//	"fmt"
	"appengine"
	"appengine/datastore"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
)

type PollApi Api

func (p PollApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path("/api/0/polls").
		Doc("Polls").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(p.getLatest).
		Doc("get the latest polls").
		Operation("getLatest").
		Writes([]db.Poll{}))
	ws.Route(ws.POST("/").To(p.createPoll).
		Filter(basicAuthenticate).
		Doc("create a poll").
		Operation("createPoll").
		Reads(db.Poll{}))

	restful.Add(ws)
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

	c.Debugf("%s", polls)
	respondMany(w, polls)
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
