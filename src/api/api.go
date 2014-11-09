package api

import (
	//	"fmt"
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"math/rand"
	"time"
)

type Api swagger.Api

type IApi interface {
	Register()
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	apis := []IApi{UserApi{}, AuthApi{}, PollApi{}, NotificationApi{}}
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
	w.WriteEntity(entities)
}

func respondError(w *restful.Response, httperr int, error string) {
	w.WriteHeader(httperr)
	w.WriteEntity(map[string]string{"error": error})
}

/*
   Auth functions
*/

func basicAuthenticate(r *restful.Request, w *restful.Response, chain *restful.FilterChain) {
	c := appengine.NewContext(r.Request)

	user, _ := auth(c, r)
	if user == nil {
		respondError(w, 401, "invalid or missing auth")
		return
	}
	c.Infof("Authenticated %s", user.Username)

	chain.ProcessFilter(r, w)
}

func auth(c appengine.Context, r *restful.Request) (*db.User, *datastore.Key) {
	authstr := r.Request.Header.Get("Authorization")
	if authstr == "" {
		authstr = r.QueryParameter("api_key")
		if authstr == "" {
			c.Warningf("got blank Authorization header and api_get GET param")
			return nil, nil
		}
	}

	var userkey *datastore.Key

	item, err := memcache.Get(c, authstr)
	if err == memcache.ErrCacheMiss {
		c.Debugf("item not in the memcache, checking in datastore")
		var auths []db.Auth
		q := datastore.NewQuery("Auth").Filter("Key =", authstr).Limit(1)
		keys, err := q.GetAll(c, &auths)
		if err != nil {
			c.Errorf(err.Error())
			panic(err)
		}
		if len(keys) == 0 {
			c.Debugf("failed to find authkey both in memcache and datastore")
			return nil, nil
		}
		key := keys[0]

		userkey = key.Parent()

		item := &memcache.Item{
			Key:   authstr,
			Value: []byte(userkey.Encode()),
		}

		if err := memcache.Set(c, item); err != nil {
			c.Errorf("error adding item: %v", err)
		}
	} else if err != nil {
		c.Errorf("error getting item: %v", err)
		return nil, nil
	} else {
		c.Debugf("found key in memcache")
		userkey, err = datastore.DecodeKey(string(item.Value))
		if err != nil {
			c.Errorf("could not decode key: %v", err)
			return nil, nil
		}
	}

	if userkey != nil {
		var user db.User

		err := datastore.Get(c, userkey, &user)
		if err != nil {
			c.Errorf("userkey %s was not in datastore", userkey)
			return nil, nil
		}
		c.Infof("User successfully authorized in with email: " + user.Email)
		return &user, userkey
	} else {
		c.Errorf("could not find Auth for authkey")
	}
	c.Errorf("something went wrong when trying to auth, authstr: %s", authstr)
	return nil, nil
}
