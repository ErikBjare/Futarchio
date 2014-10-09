package db

import (
	//	"fmt"
	"code.google.com/p/go.crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	Username string    `json:"username"`
	Password []byte    `json:"-"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
	ApiKeys  []string  `json:"-"`
}

type Auth struct {
	Key string `json:"key"`
}

func NewUser(username string, password string, name string, email string, apikeys []string) *User {
	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatal(err)
	}

	return &User{
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
