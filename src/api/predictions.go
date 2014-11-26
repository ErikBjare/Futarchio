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

// TODO: More descriptive Operation("")'s

func (p PredictionApi) Register() {
	// Statements
	ws := new(restful.WebService)
	ws.
		Path("/api/0/statements").
		Doc("Statements").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(p.getLatest).
		Doc("get the latest statements").
		Operation("getLatest").
		Writes([]StatementResponse{}))
	ws.Route(ws.POST("").To(p.createStmt).
		Doc("create a new statement").
		Operation("getLatest").
		Filter(authFilter).
		Reads(createStatement{}).
		Writes(StatementResponse{}))
	ws.Route(ws.GET("/{key}").To(p.stmtByKey).
		Doc("get the statement with given key").
		Operation("getLatest").
		Writes([]StatementResponse{}))
	ws.Route(ws.POST("/{key}/predict").To(p.predict).
		Doc("create a prediction on statement").
		Operation("getLatest").
		Writes([]StatementResponse{}))
	restful.Add(ws)

	// Predictions
	ws = new(restful.WebService)
	ws.
		Path("/api/0/predictions").
		Doc("Statements").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{key}").To(p.predByKey).
		Doc("get the prediction with given key").
		Operation("getLatest").
		Writes([]StatementResponse{}))
	restful.Add(ws)
}

func (p *PredictionApi) getLatest(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	var statements []db.Statement
	datastore.NewQuery("Statement").GetAll(c, &statements)

	respondMany(w, statements)
}

func (p *PredictionApi) createStmt(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	_, userkey := auth(c, r)

	var inputStmt createStatement
	err := r.ReadEntity(&inputStmt)
	if err != nil {
		c.Errorf("%v", err)
	}

	stmtkey := datastore.NewIncompleteKey(c, "Statement", nil)
	stmt := db.NewStatement(inputStmt.Title, inputStmt.Description, userkey)
	_, err = datastore.Put(c, stmtkey, &stmt)
	if err != nil {
		c.Errorf("%v", err)
	}

	respondOne(w, stmt)

	// TODO: Test
}

func (p *PredictionApi) stmtByKey(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	key, err := datastore.DecodeKey(r.PathParameter("key"))
	if err != nil {
		c.Errorf("%v", err)
	}

	var stmt db.Statement
	err = datastore.Get(c, key, &stmt)
	if err != nil {
		c.Errorf("%v", err)
	}

	respondOne(w, stmt)
}

func (p *PredictionApi) predByKey(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	key, err := datastore.DecodeKey(r.PathParameter("key"))
	if err != nil {
		c.Errorf("%v", err)
	}

	var pred db.Prediction
	err = datastore.Get(c, key, &pred)
	if err != nil {
		c.Errorf("%v", err)
	}

	respondOne(w, pred)
}

type credenceMsg struct {
	Credence float32 `json:"credence"`
}

func (p *PredictionApi) predict(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	stmtkey, err := datastore.DecodeKey(r.PathParameter("key"))
	if err != nil {
		c.Errorf("%v", err)
	}

	var credMsg credenceMsg
	r.ReadEntity(&credMsg)

	_, userkey := auth(c, r)

	predkey := datastore.NewIncompleteKey(c, "Predictions", nil)
	datastore.Put(c, predkey, db.NewPrediction(userkey, stmtkey, credMsg.Credence))
	// TODO: Tests
}
