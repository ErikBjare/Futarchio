package db

import (
	"code.google.com/p/go.crypto/bcrypt"
	"log"
	"strings"
	"time"
)

// A user
type User struct {
	Username string    `json:"username"`
	Password []byte    `json:"-"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
}

// Creates a new user
func NewUser(username string, password string, name string, email string) User {
	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		log.Fatal(err)
	}

	return User{
		Username: strings.ToLower(username),
		Password: hashed_pass,
		Name:     name,
		Email:    strings.ToLower(email),
		Created:  time.Now(),
	}
}

// Checks the password for a given user
func (u User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		log.Println(err)
	}
	return !(err != nil)
}
