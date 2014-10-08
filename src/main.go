package main

import (
	"github.com/ErikBjare/Futarchio/src/api"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
	"github.com/golang/oauth2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func init() {
	initUsers()
}

func initUsers() {
	for _, elem := range [][]string{{"erb", "Erik", "erik@bjareho.lt"}, {"clara", "Clara", "idunno@example.com"}} {
		username, name, email := elem[0], elem[1], elem[2]
		user := api.Users.FindOne(bson.M{"name": name})

		if user.Name != name {
			user := db.NewUser(username, "password", name, email, []string{})
			log.Println("Creating user, did not exist.\n - name: " + name + "\n - id: " + user.Id.Hex())
			err := api.Users.Insert(user)
			if err != nil {
				log.Fatal(err)
			}
		}
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
	server := &http.Server{Addr: ":80", Handler: mux}

	log.Println("Frontend is serving on: http://futarch.io/")
	log.Println("API is serving on: http://futarch.io/api/")
	server.ListenAndServe()
}

func main() {
	log.Println("Started...")
	rand.Seed(time.Now().Unix())

	wsContainer := restful.NewContainer()

	api.Users.Register(wsContainer)

	go serve(wsContainer)

	queue := make(chan error)
	for {
		err := <-queue
		log.Println(err)
	}

	log.Println("Quitting")
}
