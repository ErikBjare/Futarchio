package api

import (
	//	"fmt"
	"crypto/sha256"
	"encoding/base64"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

var (
	Users *UserResource
	Auths *AuthResource
	Polls *PollResource
)

func init() {
	session := NewSession()
	Users = NewUserResource(session)
	Auths = NewAuthResource(session)
	Polls = NewPollResource(session)
}

type MongoResource struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func NewSession() *mgo.Session {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	//defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}

func NewMongoResource(collection string, session *mgo.Session) *MongoResource {
	c := session.DB("test").C(collection)
	u := new(MongoResource)
	u.collection = c
	return u
}

/*
  UserResource
*/

type UserResource struct {
	*MongoResource
}

func NewUserResource(session *mgo.Session) *UserResource {
	ur := &UserResource{NewMongoResource("users", session)}
	return ur
}

func (ur UserResource) FindOne(bson bson.M) *db.User {
	q := ur.collection.Find(bson)
	user := db.User{}
	q.One(&user)
	return &user
}

func (ur UserResource) Insert(user *db.User) error {
	err := ur.collection.Insert(user)
	return err
}

func (ur UserResource) findById(id bson.ObjectId) *db.User {
	return ur.FindOne(bson.M{"_id": id})
}

func (u UserResource) findByAuth(authkey string) *db.User {
	auth := Auths.findByKey(authkey)
	if auth.Key != authkey {
		log.Println("Could not find matching authkey")
		return nil
	}
	user := u.findById(auth.User)
	if user.Name == "" {
		log.Println("User could not be found")
		return nil
	}
	return user
}

func (u UserResource) getByAuth(request *restful.Request, response *restful.Response) {
	authkey := request.Request.Header.Get("Authorization")
	user := u.findByAuth(authkey)
	respond(response, []db.User{*user})
}

func (u UserResource) getByKeyVal(request *restful.Request, response *restful.Response) {
	var q *mgo.Query
	key := request.PathParameter("key")
	val := request.PathParameter("val")
	if key != "" && val != "" {
		if key == "id" {
			key = "_id"
			q = u.collection.Find(bson.M{key: bson.ObjectIdHex(val)})
		} else {
			q = u.collection.Find(bson.M{key: val})
		}
	} else {
		q = u.collection.Find(bson.M{})
	}
	result := []db.User{}
	q.All(&result)
	//log.Println(fmt.Sprintf("%d matching entries in database for request: %s", len(result), request.PathParameters()))
	if len(result) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}

	respond(response, result)
}

/*
   AuthResource
*/

type AuthResource struct {
	*MongoResource
}

func NewAuthResource(session *mgo.Session) *AuthResource {
	ur := &AuthResource{NewMongoResource("auths", session)}
	return ur
}

func (ar AuthResource) Insert(auth *db.Auth) error {
	err := ar.collection.Insert(auth)
	return err
}

func (ar AuthResource) findByUserId(uid bson.ObjectId) *db.Auth {
	q := ar.collection.Find(bson.M{"user": uid})
	auth := db.Auth{}
	q.One(&auth)
	return &auth
}

func (ar AuthResource) findByKey(key string) *db.Auth {
	q := ar.collection.Find(bson.M{"key": key})
	auth := db.Auth{}
	q.One(&auth)
	return &auth
}

/*
   PollResource
*/

type PollResource struct {
	*MongoResource
}

func NewPollResource(session *mgo.Session) *PollResource {
	ur := &PollResource{NewMongoResource("polls", session)}
	return ur
}

/*
   Repspond Functions
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

func (u UserResource) Register(container *restful.Container) {
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

func (ur UserResource) authorizeUser(request *restful.Request, response *restful.Response) {
	var q *mgo.Query
	// Username can be User.Username or User.Email
	data := map[string]string{}
	err := request.ReadEntity(&data)
	if err != nil {
		log.Fatal(err)
	}
	username := data["username"]
	password := data["password"]

	if strings.Contains(username, "@") {
		// Is an Email
		q = ur.collection.Find(bson.M{"email": username})
	} else {
		q = ur.collection.Find(bson.M{"username": username})
	}

	user := db.User{}
	q.One(&user)

	if user.CheckPassword(password) {
		// If user successfully authorized

		// Check if auth key already exists
		auth := Auths.findByUserId(user.Id)
		log.Println("Found existing auth")
		log.Println(auth)

		if auth.User == user.Id {
			// If user already has a authkey
			response.WriteEntity(map[string]interface{}{"auth": auth})
		} else {
			// If user doesn't have an authkey
			auth_bytes := sha256.Sum256([]byte(username + password + strconv.Itoa(rand.Int())))
			authkey := base64.StdEncoding.EncodeToString([]byte(auth_bytes[:]))
			auth := db.Auth{User: user.Id, Key: authkey}
			log.Println("Created new auth")
			log.Println(auth)
			err := Auths.Insert(&auth)
			if err != nil {
				panic(err)
			}
			response.WriteEntity(map[string]interface{}{"auth": auth})
		}
	} else {
		response.WriteEntity(map[string]interface{}{"error": "wrong username/password"})
	}

}

func basicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	authkey := req.Request.Header.Get("Authorization")
	if len(authkey) == 0 {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}
	// usr/pwd = admin/admin
	// real code does some decoding
	auth := Auths.findByKey(authkey)
	if auth.Key == authkey {
		log.Println("Successfully logged in with UID:", auth.Id.Hex())
	} else {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}
	chain.ProcessFilter(req, resp)
}
