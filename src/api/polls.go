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
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(p.getTop).
		Filter(basicAuthenticate).
		Doc("get the latest polls").
		Operation("getTop").
		Writes([]db.Poll{}))

	restful.Add(ws)
}

func (p PollApi) getTop(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	q := datastore.NewQuery("Polls").Limit(20)

	var polls []db.Poll
	q.GetAll(c, &polls)
	respondMany(w, polls)
}
