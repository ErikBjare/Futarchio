package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	//	"strconv"
	"github.com/golang/oauth2"
	"os"
	"time"
)

type User struct {
	Id      bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Name    string        `json:"name"`
	Email   string        `json:"email"`
	ApiKeys []string      `json:"apikeys"`
}

func NewUser(name string, email string, apikeys []string) *User {
	id := bson.NewObjectId()
	return &User{Id: id, Name: name, Email: email, ApiKeys: apikeys}
}

type Resource struct {
	session    *mgo.Session
	collection *mgo.Collection
}

type UserResource struct {
	*Resource
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

func NewResource(collection string, session *mgo.Session) *Resource {
	c := session.DB("test").C(collection)
	u := new(Resource)
	u.collection = c
	return u
}

func NewUserResource(session *mgo.Session) *UserResource {
	ur := &UserResource{NewResource("users", session)}
	return ur
}

func NewPollResource(session *mgo.Session) *UserResource {
	ur := &UserResource{NewResource("users", session)}
	return ur
}

func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	var q *mgo.Query
	user_id := request.PathParameter("user-id")
	email := request.PathParameter("email")

	switch {
	case user_id != "":
		id := bson.ObjectIdHex(user_id)
		q = u.collection.Find(bson.M{"_id": id})
	case email != "":
		q = u.collection.Find(bson.M{"email": email})
	default:
		q = u.collection.Find(bson.M{})
	}

	result := []User{}
	q.All(&result)
	log.Println(fmt.Sprintf("%d matching entries in database for request: %s", len(result), request.PathParameters()))
	if len(result) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}
	response.WriteEntity(map[string]interface{}{"length": len(result), "data": result})
}

func (u UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.
		Path("/api/0/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	log.Println("Initialized paths")

	ws.Route(ws.GET("/id/{user-id}").To(u.findUser).
		Doc("get a user by id").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(User{}))
	ws.Route(ws.GET("/email/{email}").To(u.findUser).
		Doc("get a user by email").
		Param(ws.PathParameter("email", "email of the user").DataType("string")).
		Writes(User{}))
	ws.Route(ws.GET("/").To(u.findUser).
		Doc("get a list of all users").
		Writes(User{}))
	log.Println("Initialized routes")

	container.Add(ws)
}

func (u UserResource) Init() {
	c := u.collection
	for _, elem := range [][]string{{"Erik", "erik@bjareho.lt"}, {"Clara", "idunno@example.com"}} {
		name, email := elem[0], elem[1]
		result := []User{}
		err := c.Find(bson.M{"name": name}).All(&result)

		if len(result) == 0 {
			user := NewUser(name, email, []string{})
			log.Println("Creating user, did not exist.\n - name: " + name + "\n - id: " + user.Id.Hex())
			err = c.Insert(user)
			if err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Println(err)
		} else {
			log.Println(fmt.Sprintf("%d matching entries in database for name: %s, had id: %s", len(result), name, result[0].Id))
		}
	}
}

func oauth_test() {
	file, err := os.Open("key.pem")
	if err != nil {
		panic(err)
	}
	key := []byte{}
	file.Read(key)

	conf, err := oauth2.NewJWTConfig(&oauth2.JWTOptions{
		Email: "643992545442-u8dubmhq38dor5bvltjb2o98tv3musqq@developer.gserviceaccount.com",
		// The contents of your RSA private key or your PEM file
		// that contains a private key.
		// If you have a p12 file instead, you
		// can use `openssl` to export the private key into a pem file.
		//
		//    $ openssl pkcs12 -in key.p12 -out key.pem -nodes
		//
		// It only supports PEM containers with no passphrase.
		PrivateKey: key,
		Scopes:     []string{"profile"},
	},
		"https://provider.com/o/oauth2/token")
	if err != nil {
		log.Fatal(err)
	}

	// Initiate an http.Client, the following GET request will be
	// authorized and authenticated on the behalf of
	// xxx@developer.gserviceaccount.com.
	client := http.Client{Transport: conf.NewTransport()}
	client.Get("...")

	// If you would like to impersonate a user, you can
	// create a transport with a subject. The following GET
	// request will be made on the behalf of user@example.com.
	client = http.Client{Transport: conf.NewTransportWithUser("user@example.com")}
	client.Get("...")
}

func serve(wsContainer *restful.Container) {
	mux := http.NewServeMux()
	mux.Handle("/api/", wsContainer)
	mux.Handle("/", http.FileServer(http.Dir("site")))
	server := &http.Server{Addr: ":8080", Handler: mux}

	log.Println("Frontend is serving on: http://localhost:8080")
	log.Println("API is serving on: http://localhost:8080/api/")
	server.ListenAndServe()
}

func main() {
	log.Println("Started...")
	rand.Seed(time.Now().Unix())

	wsContainer := restful.NewContainer()
	session := NewSession()
	u := NewUserResource(session)
	u.Register(wsContainer)
	u.Init()

	go serve(wsContainer)

	queue := make(chan error)
	for {
		err := <-queue
		log.Println(err)
	}

	log.Println("Quitting")
}
