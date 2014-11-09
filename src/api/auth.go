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

	// Enforce lowercase to ensure case-insensitivity
	ar.Username = strings.ToLower(ar.Username)

	// Limit is two so that an error is raised upon multiple returns (which should be impossible)
	q := datastore.NewQuery("User").Limit(2)
	if strings.ContainsRune(ar.Username, '@') {
		// Is an email
		c.Debugf("Finding user to log in with email: ", ar.Username)
		q = q.Filter("Email =", ar.Username)
	} else {
		// Is a username
		c.Debugf("Finding user to log in with username: ", ar.Username)
		q = q.Filter("Username =", ar.Username)
	}

	var users []db.User
	userkeys, err := q.GetAll(c, &users)

	if len(users) > 1 {
		c.Errorf("Got more than one user when trying to auth")
		respondError(w, 500, "found more than one match for username/email, can not log in")
		return
	}

	if len(users) != 0 && users[0].CheckPassword(ar.Password) {
		// If user successfully authorized

		// Check if auth key already exists
		q := datastore.NewQuery("Auth").Ancestor(userkeys[0]).Limit(1)
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
				Value: []byte(userkeys[0].Encode()),
			}

			if err := memcache.Add(c, item); err == memcache.ErrNotStored {
				c.Infof("item with key %q already exists", item.Key)
			} else if err != nil {
				c.Errorf("error adding item: %v", err)
			}

			// Put the item into the datastore
			_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Auth", userkeys[0]), &auth)
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
