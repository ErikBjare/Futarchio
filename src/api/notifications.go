package api

import (
	"appengine"
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	"time"
)

type NotificationApi Api

func (n NotificationApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path("/api/0/notifications").
		Doc("Notifications").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(n.getNotifications).
		Doc("get the current users notifications").
		Operation("getNotifications").
		Filter(authFilter).
		Writes([]Notification{}))

	restful.Add(ws)
}

type Notification struct {
	Userkey *datastore.Key `json:"userkey"`
	Title   string         `json:"title"`
	Message string         `json:"message"`
	Created time.Time      `json:"created"`
}

func NewNotification(userkey *datastore.Key, title string, message string) Notification {
	return Notification{userkey, title, message, time.Now()}
}

func CreateNotification(c appengine.Context, user *datastore.Key, title string, message string) error {
	notif := NewNotification(user, title, message)
	key := datastore.NewIncompleteKey(c, "Notification", nil)
	key, err := datastore.Put(c, key, notif)
	if err != nil {
		return err
	}
	return nil
}

func (n *NotificationApi) getNotifications(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	_, userkey := auth(c, r)

	notifications := []Notification{}
	_, err := datastore.NewQuery("Notification").Filter("Userkey =", userkey).Limit(20).Order("-Created").GetAll(c, &notifications)
	if err != nil {
		respondError(w, 500, "error when trying to get notifications")
		return
	}

	respondMany(w, notifications)
}
