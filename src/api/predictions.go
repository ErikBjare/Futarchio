package api

import (
	"appengine"
	"appengine/datastore"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
)

type PredictionApi Api

type createStatement struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type StatementResponse struct {
	db.Statement
}

func (p PredictionApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path("/api/0/predictions").
		Doc("Predictions").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(p.getLatest).
		Doc("get the latest statements").
		Operation("getLatest").
		Writes([]StatementResponse{}))

	ws.Route(ws.POST("").To(p.create).
		Doc("create a new statement").
		Operation("getLatest").
		Filter(authFilter).
		Reads(createStatement{}).
		Writes(StatementResponse{}))

	restful.Add(ws)
}

func (p *PredictionApi) getLatest(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	var statements []db.Statement
	datastore.NewQuery("Statement").GetAll(c, &statements)

	respondMany(w, statements)
}

func (p *PredictionApi) create(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	_, userkey := auth(c, r)

	var createStmt createStatement
	err := r.ReadEntity(&createStmt)
	if err != nil {
		c.Errorf("%v", err)
	}

	stmtkey := datastore.NewIncompleteKey(c, "Statement", nil)
	stmt := db.NewCredenceStatement(createStmt.Title, createStmt.Description, userkey)
	datastore.Put(c, stmtkey, &stmt)

	respondOne(w, stmt)
}
