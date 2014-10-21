package api

import (
	//	"fmt"
	"appengine"
	"appengine/datastore"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
	"net/http"
	"strings"
)

type UserApi Api

func (u UserApi) getByAuth(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	user, _ := auth(c, r)
	if user != nil {
		respondOne(w, user)
	} else {
		w.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		w.WriteErrorString(401, "401: Not Authorized")
	}
}

func (u UserApi) getByKeyVal(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	key := r.PathParameter("key")
	val := r.PathParameter("val")

	q := datastore.NewQuery("User")
	if key != "" && val != "" {
		q = q.Filter(strings.Replace(key, key[0:1], strings.ToUpper(key[0:1]), 1)+" =", val)
	}

	result := []db.User{}
	_, err := q.GetAll(c, &result)
	if err != nil {
		panic(err)
	}

	if len(result) == 0 {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}

	respondMany(w, result)
}

func (u UserApi) create(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	var userreg UserRegistration
	err := r.ReadEntity(&userreg)
	if err != nil {
		c.Errorf(err.Error())
	}

	user := db.NewUser(userreg.Username, userreg.Password, userreg.Name, userreg.Email)
	// TODO: Write tests for username & email uniqueness
	key := datastore.NewKey(c, "User", userreg.Username, 0, nil)

	var existing_user db.User
	err = datastore.Get(c, key, &existing_user)
	if err == nil || err.Error() != "datastore: no such entity" {
		if err == nil {
			// TODO: correct errorcode
			respondError(w, 500, "username is taken")
		} else {
			respondError(w, 500, err.Error())
		}
		return
	}

	var existing_users []db.User
	keys, err := datastore.NewQuery("User").Filter("Email =", user.Email).Limit(1).GetAll(c, &existing_users)
	if err != nil {
		c.Errorf(err.Error())
		respondError(w, 500, err.Error())
		return
	} else if len(keys) != 0 {
		respondError(w, 500, "email is taken")
		return
	}

	_, err = datastore.Put(c, key, &user)
	if err != nil {
		c.Errorf(err.Error())
		respondError(w, 500, err.Error())
	}
	respondSuccess(w, "successfully created user")
}

type UserRegistration struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func (u UserApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path("/api/0/users").
		Doc("Users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(u.getByKeyVal).
		Doc("get a list of all users").
		Operation("placeholderOp").
		Filter(basicAuthenticate).
		Writes([]db.User{}))
	ws.Route(ws.POST("").To(u.create).
		Doc("create a user").
		Operation("createUser").
		Reads(UserRegistration{}))
	ws.Route(ws.GET("/me").To(u.getByAuth).
		Doc("get the authorized user").
		Operation("getByAuth").
		Filter(basicAuthenticate).
		Writes(db.User{}))
	ws.Route(ws.GET("/{key}/{val}").To(u.getByKeyVal).
		Doc("get a user by its properties").
		Operation("placeholderOp").
		Param(ws.PathParameter("key", "property to look up").DataType("string")).
		Param(ws.PathParameter("val", "value to match").DataType("string")).
		Writes([]db.User{}))

	restful.Add(ws)
}
