package api

import (
	//	"fmt"
	"appengine"
	"appengine/datastore"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
	"regexp"
	"strings"
)

type UserApi Api

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
		Writes(UserResponse{}))
	ws.Route(ws.GET("/{key}/{val}").To(u.getByKeyVal).
		Doc("get a user by its properties").
		Operation("placeholderOp").
		Param(ws.PathParameter("key", "property to look up").DataType("string")).
		Param(ws.PathParameter("val", "value to match").DataType("string")).
		Writes([]db.User{}))

	restful.Add(ws)
}

type UserResponse struct {
	db.User
	Poll_count int `json:"poll_count"`
	Vote_count int `json:"vote_count"`
}

// Serves private info such as stats
func buildUserResponse(c appengine.Context, u *db.User, key *datastore.Key) UserResponse {
	poll_count, err := datastore.NewQuery("Poll").Filter("Creator =", u.Username).Count(c)
	if err != nil {
		panic(err)
	}
	vote_count, err := datastore.NewQuery("VoteReceipt").Filter("User =", key).Count(c)
	if err != nil {
		panic(err)
	}
	return UserResponse{*u, poll_count, vote_count}
}

func (u UserApi) getByAuth(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	user, key := auth(c, r)
	if user != nil {
		ur := buildUserResponse(c, user, key)
		respondOne(w, ur)
	} else {
		respondError(w, 500, "something went wrong when trying to get user")
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
		respondError(w, 404, "user could not be found")
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

	// TODO: Validate email properly
	if userreg.Email == "" {
		respondError(w, 500, "email field was empty")
		return
	}
	if !strings.Contains(userreg.Email, "@") {
		respondError(w, 500, "email did not contain a '@'")
	}

	// TODO: Move username validation to seperate function ValidateUsername (probably in db)
	matched, err := regexp.MatchString("^[a-z0-9_]{3,20}$", userreg.Username)
	if err != nil {
		panic(err)
	}
	if !matched {
		respondError(w, 500, "username contains invalid characters, can only contain 3-20 lowercase a-z, 0-9 and _")
		return
	}

	user := db.NewUser(userreg.Username, userreg.Password, userreg.Name, userreg.Email)
	// TODO: Write tests for username & email uniqueness, also make sure they match their regexps
	key := datastore.NewKey(c, "User", user.Username, 0, nil)

	var existing_user db.User
	err = datastore.Get(c, key, &existing_user)
	if err == nil {
		respondError(w, 500, "username is taken")
		return
	}
	if err.Error() != "datastore: no such entity" {
		respondError(w, 500, err.Error())
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
