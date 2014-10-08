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
)

func init() {
	session := NewSession()
	Auths = &AuthResource{map[string]bson.ObjectId{}}
	Users = NewUserResource(session)
}

type MongoResource struct {
	session    *mgo.Session
	collection *mgo.Collection
}

type UserResource struct {
	*MongoResource
}

type AuthResource struct {
	auths map[string]bson.ObjectId
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

func NewUserResource(session *mgo.Session) *UserResource {
	ur := &UserResource{NewMongoResource("users", session)}
	return ur
}

func NewPollResource(session *mgo.Session) *UserResource {
	ur := &UserResource{NewMongoResource("users", session)}
	return ur
}

func (ur UserResource) FindUserById(id bson.ObjectId) db.User {
	q := ur.collection.Find(bson.M{"_id": id})
	user := db.User{}
	q.One(&user)
	return user
}

func (u UserResource) byAuth(request *restful.Request, response *restful.Response) {
	var q *mgo.Query
	auth := request.Request.Header.Get("Authorization")
	uid, ok := Auths.auths[auth]
	if ok {
		q = u.collection.Find(bson.M{"_id": uid})
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

func (u UserResource) byKeyVal(request *restful.Request, response *restful.Response) {
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

func respond(r *restful.Response, result interface{}) {
	r.WriteEntity(map[string]interface{}{"data": result})
}

func (u UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.
		Path("/api/0/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/me").To(u.byAuth).
		Filter(basicAuthenticate).
		Doc("get the authorized user").
		Writes(db.User{}))
	ws.Route(ws.GET("/{key}/{val}").To(u.byKeyVal).
		Doc("get a user by its properties").
		Param(ws.PathParameter("key", "property to look up").DataType("string")).
		Param(ws.PathParameter("val", "value to match").DataType("string")).
		Writes(db.User{}))
	ws.Route(ws.GET("/").To(u.byKeyVal).
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

func authFromUserId(uid bson.ObjectId) string {
	for k, v := range Auths.auths {
		if v == uid {
			return k
		}
	}
	return ""
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
		auth := authFromUserId(user.Id)

		if auth == "" {
			auth_bytes := sha256.Sum256([]byte(username + password + strconv.Itoa(rand.Int())))
			auth = base64.StdEncoding.EncodeToString([]byte(auth_bytes[:]))
			Auths.auths[auth] = user.Id
			response.WriteEntity(map[string]interface{}{"auth": auth})
		} else {
			response.WriteEntity(map[string]interface{}{"auth": auth})
		}
	} else {
		response.WriteEntity(map[string]interface{}{"error": "wrong username/password"})
	}

}

func basicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	encoded := req.Request.Header.Get("Authorization")
	// usr/pwd = admin/admin
	// real code does some decoding
	uid, ok := Auths.auths[encoded]
	if len(encoded) == 0 {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	} else if "Basic YWRtaW46YWRtaW4=" == encoded {
		log.Print("Successfully logged in with admin,admin")
	} else if ok {
		log.Println("Successfully logged in with UID:", uid.Hex())
	} else {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}
	chain.ProcessFilter(req, resp)
}

func (u UserResource) Init() {
	c := u.collection
	for _, elem := range [][]string{{"erb", "Erik", "erik@bjareho.lt"}, {"clara", "Clara", "idunno@example.com"}} {
		username, name, email := elem[0], elem[1], elem[2]
		result := []db.User{}
		err := c.Find(bson.M{"name": name}).All(&result)

		if len(result) == 0 {
			user := db.NewUser(username, "password", name, email, []string{})
			log.Println("Creating user, did not exist.\n - name: " + name + "\n - id: " + user.Id.Hex())
			err = c.Insert(user)
			if err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Println(err)
		} /* else {
			log.Println(fmt.Sprintf("%d matching entries in database for name: %s, had id: %s", len(result), name, result[0].Id))
		}*/
	}
}
