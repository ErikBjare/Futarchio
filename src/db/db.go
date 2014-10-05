package db

import (
	//	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	// "strings"
	//	"strconv"
	"code.google.com/p/go.crypto/bcrypt"
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
