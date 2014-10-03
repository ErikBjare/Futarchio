package main

import (
	//	"fmt"
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"strings"
	//	"strconv"
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/golang/oauth2"
	"os"
	"time"
)

type User struct {
	Id       bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Username string        `json:"username"`
	Password []byte        `json:"password"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Created  time.Time     `json:"created"`
	ApiKeys  []string      `json:"apikeys"`
}

func NewUser(username string, password string, name string, email string, apikeys []string) *User {
	id := bson.NewObjectId()
	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatal(err)
	}

	return &User{
		Id:       id,
		Username: username,
		Password: hashed_pass,
		Name:     name,
		Email:    email,
		Created:  time.Now(),
		ApiKeys:  apikeys,
	}
}

func (u User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		log.Println(err)
	}
	return !(err != nil)
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

	result := []User{}
	q.All(&result)
	//log.Println(fmt.Sprintf("%d matching entries in database for request: %s", len(result), request.PathParameters()))
	if len(result) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}

	response.WriteEntity(map[string]interface{}{"length": len(result), "data": result})
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
	log.Println(data)

	if strings.Contains(username, "@") {
		// Is an Email
		q = ur.collection.Find(bson.M{"email": username})
	} else {
		q = ur.collection.Find(bson.M{"username": username})
	}

	user := User{}
	q.One(&user)

	if user.CheckPassword(string(password)) {
		// If user successfully authorized
		response.WriteEntity(map[string]interface{}{"data": user})
	} else {
		response.WriteEntity(map[string]interface{}{"error": "wrong username/password"})
	}

}

func basicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	encoded := req.Request.Header.Get("Authorization")
	// usr/pwd = admin/admin
	// real code does some decoding
	if len(encoded) == 0 || "Basic YWRtaW46YWRtaW4=" != encoded {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}
	chain.ProcessFilter(req, resp)
}

func (u UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.
		Path("/api/0/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{key}/{val}").To(u.findUser).
		Doc("get a user by its properties").
		Param(ws.PathParameter("key", "property to look up").DataType("string")).
		Param(ws.PathParameter("val", "value to match").DataType("string")).
		Writes(User{}))
	ws.Route(ws.GET("/").To(u.findUser).
		Filter(basicAuthenticate).
		Doc("get a list of all users").
		Writes(User{}))
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

func (u UserResource) Init() {
	c := u.collection
	for _, elem := range [][]string{{"erb", "Erik", "erik@bjareho.lt"}, {"clara", "Clara", "idunno@example.com"}} {
		username, name, email := elem[0], elem[1], elem[2]
		result := []User{}
		err := c.Find(bson.M{"name": name}).All(&result)

		if len(result) == 0 {
			user := NewUser(username, "password", name, email, []string{})
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

func oauth_test() {
	file, err := os.Open("secrets/key.pem")
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
	mux.Handle("/api/0/", wsContainer)
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
