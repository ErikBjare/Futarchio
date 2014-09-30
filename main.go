package main

import (
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	_Id  bson.ObjectId `json:"id"`
	Name string        `json:"name"`
}

func (u User) getJson() map[string]string {
	return map[string]string{"name": u.Name, "id": strconv.Itoa(u.Id)}
}

type Resource struct {
	session    *mgo.Session
	collection *mgo.Collection
}

type UserResource struct {
	Resource
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

func NewResource(collection string, session *mgo.Session) *UserResource {
	c := session.DB("test").C(collection)
	u := new(UserResource)
	u.collection = c
	return u
}

func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	id, err := strconv.Atoi(request.PathParameter("user-id"))
	if err != nil {
		log.Fatal(err)
	} else {
		result := []User{}
		err = u.collection.Find(bson.M{"id": id}).All(&result)
		log.Println(result[0].Name)
		response.WriteAsJson(result)
	}
}

func paths(ws *restful.WebService) {
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	log.Println("Initialized paths")
}

func routes(ws *restful.WebService) {
	session := NewSession()
	u := NewResource("", session)
	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(User{}))
	log.Println("Initialized routes")
}

func serve(ws *restful.WebService) {
	container := restful.NewContainer()
	container.Add(ws)
	server := &http.Server{Addr: ":8081", Handler: container}
	log.Println("Listening...")
	server.ListenAndServe()
}

func test_mgo() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	for _, element := range []string{"Erik", "Clara"} {
		c := session.DB("test").C("users")

		result := User{}
		err = c.Find(bson.M{"name": element}).One(&result)

		if result == *new(User) {
			id := rand.Int()
			log.Println("Creating user, did not exist.\n - name: " + element + "\n - id: " + strconv.Itoa(id))
			err = c.Insert(&User{id, element})
			if err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Println(err)
		} else {
			log.Println("In database: " + result.Name)
		}

	}
}

func main() {
	log.Println("Started...")
	rand.Seed(time.Now().Unix())
	ws := new(restful.WebService)

	test_mgo()

	paths(ws)
	routes(ws)
	serve(ws)

	log.Println("Quitting")
}
