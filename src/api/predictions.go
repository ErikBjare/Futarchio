package api

import (
	"appengine"
	"appengine/datastore"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
)

type PredictionApi Api

type StatementCreator struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type StatementResponse struct {
	db.Statement
	Key *datastore.Key `json:"key"`
}

type PredictionCreator struct {
	Credence float32 `json:"credence"`
}

type PredictionResponse struct {
	db.Prediction
	Key *datastore.Key `json:"key"`
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
		Operation("latest").
		Writes([]StatementResponse{}))
	ws.Route(ws.POST("").To(p.createStmt).
		Doc("create a new statement").
		Operation("create").
		Filter(authFilter).
		Reads(StatementCreator{}).
		Writes(StatementResponse{}))
	ws.Route(ws.GET("/{key}").To(p.stmtByKey).
		Doc("get the statement with given key").
		Operation("getByKey").
		Param(ws.PathParameter("key", "key of statement to fetch").DataType("string")).
		Writes(StatementResponse{}))
	ws.Route(ws.GET("/{key}/predictions").To(p.predsByStmt).
		Doc("get the statement with given key").
		Operation("getByKey").
		Param(ws.PathParameter("key", "key of statement to fetch").DataType("string")).
		Writes([]PredictionResponse{}))
	ws.Route(ws.POST("/{key}/predict").To(p.predict).
		Doc("create a prediction on statement").
		Operation("predict").
		Param(ws.PathParameter("key", "key of statement to fetch").DataType("string")).
		Reads(PredictionCreator{}).
		Writes(PredictionResponse{}))
	restful.Add(ws)

	// Predictions
	ws = new(restful.WebService)
	ws.
		Path("/api/0/predictions").
		Doc("Predictions").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{key}").To(p.predByKey).
		Doc("get the prediction with given key").
		Operation("getByKey").
		Writes(PredictionResponse{}))
	restful.Add(ws)
}

func (p *PredictionApi) getLatest(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	var statements []db.Statement
	keys, err := datastore.NewQuery("Statement").Order("-Created").GetAll(c, &statements)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}

	statementresps := make([]StatementResponse, len(statements))
	for i := range keys {
		statementresps[i].Statement = statements[i]
		statementresps[i].Key = keys[i]
	}

	respondMany(w, statementresps)
}

func (p *PredictionApi) createStmt(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	_, userkey := auth(c, r)

	var inputStmt StatementCreator
	err := r.ReadEntity(&inputStmt)
	if err != nil {
		c.Errorf("%v", err)
	}

	key := datastore.NewIncompleteKey(c, "Statement", nil)
	stmt := db.NewStatement(inputStmt.Title, inputStmt.Description, userkey)
	key, err = datastore.Put(c, key, &stmt)
	if err != nil {
		c.Errorf("%v", err)
	}

	respondOne(w, StatementResponse{stmt, key})
	// TODO: Test
}

func (p *PredictionApi) predsByStmt(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	key, err := datastore.DecodeKey(r.PathParameter("key"))
	if err != nil {
		c.Errorf("%v", err)
	}

	var preds []db.Prediction
	keys, err := datastore.NewQuery("Prediction").Filter("Statement =", key).Order("-Created").GetAll(c, &preds)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}

	predsresp := make([]PredictionResponse, len(preds))
	for i := range preds {
		predsresp[i].Prediction = preds[i]
		predsresp[i].Key = keys[i]
	}

	respondMany(w, predsresp)
}

func (p *PredictionApi) stmtByKey(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	key, err := datastore.DecodeKey(r.PathParameter("key"))
	if err != nil {
		c.Errorf("%v", err)
	}

	var stmt StatementResponse
	stmt.Key = key
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

	var pred PredictionResponse
	pred.Key = key
	err = datastore.Get(c, key, &pred)
	if err != nil {
		c.Errorf("%v", err)
	}

	respondOne(w, pred)
}

func (p *PredictionApi) predict(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	_, userkey := auth(c, r)

	// TODO: Check so that the user hasn't already made a prediction (or should historical predictions be preserved? Probably.)

	stmtkey, err := datastore.DecodeKey(r.PathParameter("key"))
	if err != nil {
		c.Errorf("%v", err)
	}

	var predcreator PredictionCreator
	err = r.ReadEntity(&predcreator)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}

	// Check if credence in valid range
	if !(0 < predcreator.Credence && predcreator.Credence < 1) {
		respondError(w, 500, "credence was not in valid range")
		return
	}

	predkey := datastore.NewIncompleteKey(c, "Prediction", nil)
	pred := db.NewPrediction(userkey, stmtkey, predcreator.Credence)
	predkey, err = datastore.Put(c, predkey, &pred)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}

	predresp := PredictionResponse{pred, predkey}
	respondOne(w, predresp)
	// TODO: Tests
}
