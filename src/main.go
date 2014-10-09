package main

import (
	"code.google.com/p/gcfg"
	"github.com/ErikBjare/Futarchio/src/api"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
	"github.com/golang/oauth2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
)

type config struct {
	Main struct {
		Hostname string
		Port     string
	}
}

var (
	Config config
)

func init() {
	err := gcfg.ReadFileInto(&Config, "config.ini")
	if err != nil {
		panic(err)
	}

	initUsers()

	serve()
}

func serve() {
	wsContainer := restful.NewContainer()
	api.Users.Register(wsContainer)
	http.Handle("/api/0/", wsContainer)

	log.Println("Frontend is serving on: http://localhost:" + Config.Main.Port)
	log.Println("API is serving on: http://localhost:" + Config.Main.Port + "/api/")
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
