package frontend

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	_ "github.com/ErikBjare/Futarchio/src/api"
	"github.com/ErikBjare/Futarchio/src/db"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"log"
	"net/http"
)

func main() {
	log.Println("Invalid use, start with `goapp serve` (Google App Engine)")
}

func init() {
	serve()
}

func getGaeURL() string {
	if appengine.IsDevAppServer() {
		return "http://localhost:8080"
	} else {
		/**
		 * Include your URL on App Engine here.
		 * I found no way to get AppID without appengine.Context and this always
		 * based on a http.Request.
		 */
		return "http://futarchio.appspot.com"
	}
}

func serve() {
	log.Println("Instance starting...")
	http.HandleFunc("/api/0/init", initDB)

	config := swagger.Config{
		WebServices:    restful.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: getGaeURL(),
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath: "/apidocs/",
		// GAE support static content which is configured in your app.yaml.
		// This example expect the swagger-ui in static/swagger so you should place it there :)
		SwaggerFilePath: "site/swagger/dist"}
	swagger.InstallSwaggerService(config)
}

func initDB(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, elem := range [][]string{{"erb", "Erik", "erik@bjareho.lt"}, {"clara", "Clara", "idunno@example.com"}} {
		username, name, email := elem[0], elem[1], elem[2]
		q := datastore.NewQuery("User").
			Filter("email =", email).
			Limit(1)

		var users []db.User
		q.GetAll(c, &users)

		if len(users) == 0 {
			user := db.NewUser(username, "password", name, email)
			log.Println("Creating user, did not exist.\n - name: " + name)
			key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "User", nil), user)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "created %q\n", key)
		}
	}
}
