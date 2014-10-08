package db

import (
	//	"fmt"
	"code.google.com/p/go.crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
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

type User struct {
	Id       bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Username string        `json:"username"`
	Password []byte        `json:"-"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Created  time.Time     `json:"created"`
	ApiKeys  []string      `json:"-"`
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
