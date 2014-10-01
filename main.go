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
	"time"
)

type User struct {
	Id   bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Name string        `json:"name"`
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
	id := bson.ObjectIdHex(request.PathParameter("user-id"))
	result := []User{}
	u.collection.Find(bson.M{"_id": id}).All(&result)
	log.Println(fmt.Sprintf("%d matching entries in database for id: %s", len(result), id))
	response.WriteAsJson(map[string]interface{}{"length": len(result), "data": result})
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
	u := NewResource("users", session)
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
	session := NewSession()
	userResource := NewResource("users", session)
	c := userResource.collection
	for _, element := range []string{"Erik", "Clara"} {
		result := []User{}
		err := c.Find(bson.M{"name": element}).All(&result)

		if len(result) == 0 {
			id := bson.NewObjectId()
			log.Println("Creating user, did not exist.\n - name: " + element + "\n - id: " + id.Hex())
			err = c.Insert(&User{id, element})
			if err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Println(err)
		} else {
			log.Println(fmt.Sprintf("%d matching entries in database for name: %s, had id: %s", len(result), element, result[0].Id))
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
