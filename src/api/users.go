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
	ws.Route(ws.GET("").To(u.query).
		Doc("get a list of all users").
		Operation("getAll").
		Param(ws.QueryParameter("key", "key of user to get").DataType("string")).
		Param(ws.QueryParameter("username", "username of user to get").DataType("string")).
		Writes([]db.User{}))
	ws.Route(ws.POST("").To(u.create).
		Doc("create a user").
		Operation("createUser").
		Reads(UserRegistration{}))
	ws.Route(ws.GET("/me").To(u.getByAuth).
		Doc("get the authorized user").
		Operation("getByAuth").
		Filter(authFilter).
		Writes(UserResponse{}))

	restful.Add(ws)
}

type UserRegistration struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type UserResponse struct {
	db.User
	Poll_count int            `json:"poll_count"`
	Vote_count int            `json:"vote_count"`
	Key        *datastore.Key `json:"key"`
}

// Serves private info such as stats
func buildUserResponse(c appengine.Context, u *db.User, key *datastore.Key) UserResponse {
	poll_count, err := datastore.NewQuery("Poll").Filter("Creator =", key).Count(c)
	if err != nil {
		panic(err)
	}
	vote_count, err := datastore.NewQuery("VoteReceipt").Filter("User =", key).Count(c)
	if err != nil {
		panic(err)
	}
	return UserResponse{*u, poll_count, vote_count, key}
}

func (u *UserApi) getByAuth(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	user, key := auth(c, r)
	if user != nil {
		ur := buildUserResponse(c, user, key)
		respondOne(w, ur)
	} else {
		respondError(w, 500, "something went wrong when trying to get user")
	}
}

func (u *UserApi) query(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	keystr := r.QueryParameter("key")
	username := r.QueryParameter("username")

	var users []db.User
	if keystr != "" {
		key, err := datastore.DecodeKey(keystr)
		if err != nil {
			respondError(w, 500, err.Error())
			return
		}

		var user db.User
		err = datastore.Get(c, key, &user)
		if err != nil {
			respondError(w, 500, err.Error())
			return
		}

		users = []db.User{user}
	} else {
		q := datastore.NewQuery("User")
		if username != "" {
			q = q.Filter("Username =", username)
		}
		_, err := q.GetAll(c, &users)
		if err != nil {
			panic(err)
		}
	}

	if len(users) == 0 {
		respondError(w, 404, "user could not be found")
		return
	}

	respondMany(w, users)
}

func (u *UserApi) create(r *restful.Request, w *restful.Response) {
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
		return
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
