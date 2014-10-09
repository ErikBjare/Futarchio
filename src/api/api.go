package api

import (
	//	"fmt"
	"appengine"
	"appengine/datastore"
	"crypto/sha256"
	"encoding/base64"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

var (
	Users *UserApi
	Auths *AuthApi
	Polls *PollApi
)

type UserApi struct {
	Path string
}

func init() {
	Users = &UserApi{Path: "/users"}
	Auths = &AuthApi{Path: "/auth"}
	Polls = &PollApi{Path: "/polls"}
}

func (u UserApi) getByAuth(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	authkey := r.Request.Header.Get("Authorization")

	q := datastore.NewQuery("Auth").Filter("key =", authkey).Limit(1)
	var auths []db.Auth
	keys, err := q.GetAll(c, &auths)
	if err != nil {
		panic(err)
	}

	if len(auths) != 0 && auths[0].Key == authkey {
		var user db.User
		err := datastore.Get(c, keys[0].Parent(), &user)
		if err != nil {
			panic(err)
		}
		log.Println("Successfully logged in with email: " + user.Email)
		respond(w, []db.User{user})
	} else {
		w.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		w.WriteErrorString(401, "401: Not Authorized")
		return
	}
}

func (u UserApi) getByKeyVal(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	key := r.PathParameter("key")
	val := r.PathParameter("val")

	q := datastore.NewQuery("User")
	if key != "" && val != "" {
		q = q.Filter(key+" =", val)
	} else {
		// TODO: Get all
	}
	result := []db.User{}
	_, err := q.GetAll(c, &result)
	if err != nil {
		panic(err)
	}
	//log.Println(fmt.Sprintf("%d matching entries in database for r: %s", len(result), r.PathParameters()))
	if len(result) == 0 {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}

	respond(w, result)
}

/*
   AuthApi
*/

type AuthApi struct {
	Path string
}

/*
   PollApi
*/

type PollApi struct {
	Path string
}

/*
   Respond Functions
*/

func respond(r *restful.Response, result interface{}) {
	r.WriteEntity(map[string]interface{}{"data": result})
}

func respondError(r *restful.Response, error string) {
	r.WriteHeader(http.StatusNotFound)
	r.WriteEntity(map[string]interface{}{"error": error})
}

/*
   Register functions
*/

func (u UserApi) Register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.
		Path("/api/0/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/me").To(u.getByAuth).
		Filter(basicAuthenticate).
		Doc("get the authorized user").
		Writes(db.User{}))
	ws.Route(ws.GET("/{key}/{val}").To(u.getByKeyVal).
		Doc("get a user by its properties").
		Param(ws.PathParameter("key", "property to look up").DataType("string")).
		Param(ws.PathParameter("val", "value to match").DataType("string")).
		Writes(db.User{}))
	ws.Route(ws.GET("/").To(u.getByKeyVal).
		Filter(basicAuthenticate).
		Doc("get a list of all users").
		Writes(db.User{}))
	container.Add(ws)

	ws = new(restful.WebService)
	ws.
		Path("/api/0/auth").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.POST("").To(u.authorizeUser).
		Doc("authorize a user").
		Reads(map[string]string{}))
	container.Add(ws)

	log.Println("Initialized routes")
	log.Println("Initialized paths")
}

func (ur UserApi) authorizeUser(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	// Username can be User.Username or User.Email
	data := map[string]string{}
	err := r.ReadEntity(&data)
	if err != nil {
		log.Fatal(err)
	}
	username := data["username"]
	password := data["password"]

	q := datastore.NewQuery("User").Limit(1)
	if strings.Contains(username, "@") {
		// Is an Email
		q = q.Filter("email =", username)
	} else {
		q = q.Filter("username =", username)
	}

	var user []db.User
	userkey, err := q.GetAll(c, &user)

	if len(user) != 0 && user[0].CheckPassword(password) {
		// If user successfully authorized

		// Check if auth key already exists
		q := datastore.NewQuery("Auth").Ancestor(userkey[0]).Limit(1)
		var auths []db.Auth
		_, err := q.GetAll(c, &auths)
		if err != nil {
			panic(err)
		}

		if len(auths) != 0 {
			// If user already has a authkey
			log.Println("Found existing auth")
			log.Println(auths[0])
			w.WriteEntity(map[string]interface{}{"auth": auths[0]})
		} else {
			// If user doesn't have an authkey
			auth_bytes := sha256.Sum256([]byte(username + password + strconv.Itoa(rand.Int())))
			authkey := base64.StdEncoding.EncodeToString([]byte(auth_bytes[:]))
			auth := db.Auth{Key: authkey}
			log.Println("Created new auth")
			log.Println(auth)
			_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Auth", userkey[0]), &auth)
			if err != nil {
				panic(err)
			}
			w.WriteEntity(map[string]interface{}{"auth": auth})
		}
	} else {
		w.WriteEntity(map[string]interface{}{"error": "wrong username/password"})
	}

}

func basicAuthenticate(r *restful.Request, w *restful.Response, chain *restful.FilterChain) {
	c := appengine.NewContext(r.Request)
	authkey := r.Request.Header.Get("Authorization")
	if len(authkey) == 0 {
		w.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		w.WriteErrorString(401, "401: Not Authorized")
		return
	}
	// usr/pwd = admin/admin
	// real code does some decoding
	q := datastore.NewQuery("Auth").Filter("key =", authkey).Limit(1)
	var auths []db.Auth
	keys, err := q.GetAll(c, &auths)
	if err != nil {
		panic(err)
	}

	if len(auths) != 0 && auths[0].Key == authkey {
		var user db.User
		err := datastore.Get(c, keys[0].Parent(), &user)
		if err != nil {
			panic(err)
		}
		log.Println("Successfully logged in with email: " + user.Email)
	} else {
		w.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		w.WriteErrorString(401, "401: Not Authorized")
		return
	}
	chain.ProcessFilter(r, w)
}
