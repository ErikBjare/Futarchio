package api

import (
	//	"fmt"
	"appengine"
	"appengine/datastore"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

var (
	Users *UserApi
	Auths *AuthApi
	Polls *PollApi
)

func init() {
	Users = &UserApi{}
	Auths = &AuthApi{}
	Polls = &PollApi{}

	Users.Register()
	Auths.Register()
	Polls.Register()
}

type Api swagger.Api

/*
   Respond Functions
*/

func respondOne(w *restful.Response, entity interface{}) {
	w.WriteEntity(entity)
}

func respondMany(w *restful.Response, entities interface{}) {
	w.WriteEntity(entities)
}

func respondError(w *restful.Response, httperr int, error string) {
	w.WriteHeader(httperr)
	w.WriteEntity(map[string]interface{}{"error": error})
}

/*
   Auth functions
*/

func basicAuthenticate(r *restful.Request, w *restful.Response, chain *restful.FilterChain) {
	c := appengine.NewContext(r.Request)
	authkey := r.Request.Header.Get("Authorization")

	user := auth(c, authkey)
	if user == nil {
		respondError(w, 401, "invalid or missing Authorization header")
		return
	}
	c.Infof("Authenticated %s", user.Username)

	chain.ProcessFilter(r, w)
}

func auth(c appengine.Context, authkey string) *db.User {
	q := datastore.NewQuery("Auth").Filter("Key =", authkey).Limit(1)
	var auths []db.Auth

	// TODO: Complete the following memcache implementation
	/*
		item, err := memcache.Get(c, authkey)
		if err == memcache.ErrCacheMiss {
			c.Infof("item not in the cache")
		} else if err != nil {
			c.Errorf("error getting item: %v", err)
		}

		var bytes Buffer
		var userkey datastore.Key
		dec := gob.NewDecoder(&bytes)
		err = dec.Decode(item.Key)
		if err != nil {
			c.Errorf(err)
		}
	*/

	keys, err := q.GetAll(c, &auths)
	if err != nil {
		c.Errorf(err.Error())
		panic(err)
	}
	if len(keys) == 0 {
		c.Infof("Failed to find authkey")
		return nil
	}
	key := keys[0]

	auth := db.Auth{}
	err = datastore.Get(c, key, &auth)
	if err != nil {
		c.Errorf(err.Error())
	}

	if auth.Key == authkey {
		var user db.User

		err := datastore.Get(c, key.Parent(), &user)
		if err != nil {
			panic(err)
		}
		c.Infof("User successfully authorized in with email: " + user.Email)
		return &user
	} else {
		c.Warningf("Could not find Auth for authkey")
	}
	return nil
}
