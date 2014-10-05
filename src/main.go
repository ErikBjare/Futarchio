package main

import (
	"./db"
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/base64"
	"github.com/emicklei/go-restful"
	"github.com/golang/oauth2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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

var (
	userresource *UserResource
	authresource *AuthResource
)

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

	result := []db.User{}
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
		auth := ""
		for k, v := range authresource.auths {
			if v == user.Id {
				auth = k
				break
			}
		}

		if auth == "" {
			auth_bytes, err := bcrypt.GenerateFromPassword([]byte(username+password+strconv.Itoa(rand.Int())), 8)
			if err != nil {
				log.Println(err)
			}
			auth = base64.StdEncoding.EncodeToString(auth_bytes)
			authresource.auths[auth] = user.Id
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
	uid, ok := authresource.auths[encoded]
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
		Writes(db.User{}))
	ws.Route(ws.GET("/").To(u.findUser).
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
	mux.Handle("/", http.FileServer(http.Dir("site/dist")))
	server := &http.Server{Addr: ":8080", Handler: mux}

	log.Println("Frontend is serving on: http://localhost:8080")
	log.Println("API is serving on: http://localhost:8080/api/")
	server.ListenAndServe()
}

func FindUserById(id bson.ObjectId) db.User {
	ur := userresource
	q := ur.collection.Find(bson.M{"_id": id})
	user := db.User{}
	q.One(&user)
	return user
}

func main() {
	log.Println("Started...")
	rand.Seed(time.Now().Unix())

	session := NewSession()
	wsContainer := restful.NewContainer()

	authresource = &AuthResource{map[string]bson.ObjectId{}}
	userresource = NewUserResource(session)
	userresource.Register(wsContainer)
	userresource.Init()

	go serve(wsContainer)

	queue := make(chan error)
	for {
		err := <-queue
		log.Println(err)
	}

	log.Println("Quitting")
}
