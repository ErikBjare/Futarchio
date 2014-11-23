package api

import (
	"appengine"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"math/rand"
	"reflect"
	"time"
)

type Api swagger.Api

type IApi interface {
	Register()
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	apis := []IApi{UserApi{}, AuthApi{}, PollApi{}, NotificationApi{}, PredictionApi{}}
	for _, api := range apis {
		api.Register()
	}
}

/*
   Respond Functions
*/

func respondOne(w *restful.Response, entity interface{}) {
	w.WriteEntity(entity)
}

func respondSuccess(w *restful.Response, msg string) {
	w.WriteHeader(200)
	w.WriteEntity(map[string]interface{}{"success": msg})
}

func respondMany(w *restful.Response, entities interface{}) {
	if reflect.ValueOf(entities).Len() == 0 {
		entities = []string{}
	}
	w.WriteEntity(entities)
}

func respondError(w *restful.Response, httperr int, error string) {
	w.WriteHeader(httperr)
	w.WriteEntity(map[string]string{"error": error})
}

/*
   Auth functions
*/

func authFilter(r *restful.Request, w *restful.Response, chain *restful.FilterChain) {
	c := appengine.NewContext(r.Request)

	user, _ := auth(c, r)
	if user == nil {
		respondError(w, 401, "invalid or missing auth")
		return
	}
	c.Infof("Authenticated %s", user.Username)

	chain.ProcessFilter(r, w)
}

func authAdminFilter(r *restful.Request, w *restful.Response, chain *restful.FilterChain) {
	c := appengine.NewContext(r.Request)

	user, _ := auth(c, r)
	if user == nil {
		respondError(w, 401, "invalid or missing auth")
		return
	}

	if !user.Admin() {
		respondError(w, 401, "you lack the required priviledges")
		return
	}

	c.Infof("Authenticated admin %s", user.Username)

	chain.ProcessFilter(r, w)
}
