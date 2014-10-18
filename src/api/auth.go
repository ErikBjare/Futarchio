package api

import (
	//	"fmt"
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
	"strings"
)

type AuthApi Api

type AuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a AuthApi) Register() {
	ws := new(restful.WebService)
	ws.
		Path("/api/0/auth").
		Doc("Authentication Tokens").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.POST("/").To(a.authorizeUser).
		Doc("authorize a user").
		Operation("authorizeUser").
		Reads(AuthReq{}))

	restful.Add(ws)
}

func (a AuthApi) authorizeUser(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	// Username can be User.Username or User.Email
	var ar AuthReq
	err := r.ReadEntity(&ar)
	if err != nil {
		c.Criticalf(err.Error())
		panic(err)
	}

	q := datastore.NewQuery("User").Limit(1)
	if strings.Contains(ar.Username, "@") {
		// Is an Email
		q = q.Filter("Email =", ar.Username)
	} else {
		q = q.Filter("Username =", ar.Username)
	}

	var user []db.User
	userkey, err := q.GetAll(c, &user)

	if len(user) != 0 && user[0].CheckPassword(ar.Password) {
		// If user successfully authorized

		// Check if auth key already exists
		q := datastore.NewQuery("Auth").Ancestor(userkey[0]).Limit(1).EventualConsistency()
		var auths []db.Auth
		k, err := q.GetAll(c, &auths)
		if len(k) != 0 && err != nil {
			panic(err)
		}

		if len(auths) != 0 {
			// If user already has a authkey
			respondOne(w, auths[0])
		} else {
			// If user doesn't have an authkey
			auth := db.NewAuth()

			// Add it to the memcache
			item := &memcache.Item{
				Key:   auth.Key,
				Value: []byte(userkey[0].Encode()),
			}

			if err := memcache.Add(c, item); err == memcache.ErrNotStored {
				c.Infof("item with key %q already exists", item.Key)
			} else if err != nil {
				c.Errorf("error adding item: %v", err)
			}

			// Put the item into the datastore
			_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Auth", userkey[0]), &auth)
			if err != nil {
				panic(err)
			}
			c.Infof("Created new auth %s", auth)

			respondOne(w, auth)
		}
	} else {
		respondError(w, 401, "wrong username/password")
	}

}
