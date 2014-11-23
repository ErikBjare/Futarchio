package db

import (
	"code.google.com/p/go.crypto/bcrypt"
	"log"
	"strings"
	"time"
)

// User - A user
type User struct {
	Username string    `json:"username"`
	Password []byte    `json:"-"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
	// Rank should never be accessed directly, always use getters and setters.
	Rank int `json:"-"`
}

// NewUser - Creates a new user
func NewUser(username string, password string, name string, email string) User {
	user := User{
		Username: strings.ToLower(username),
		Name:     name,
		Email:    strings.ToLower(email),
		Created:  time.Now(),
	}

	err := user.SetPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	return user
}

// SetPassword - Sets a new password for the given user
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	u.Password = hashedPassword
	return err
}

// CheckPassword - Checks the password for a given user
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		log.Println(err)
	}
	return !(err != nil)
}

// Admin - Returns true if user is admin, else false.
func (u *User) Admin() bool {
	return u.Rank >= 1000
}
