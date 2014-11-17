package db

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"
)

// Auth - Represents an authentication key
//
// Should always be an ancestor of User in the datastore
type Auth struct {
	// Should always be an ancestor of User
	Key     string    `json:"key"`
	Expires time.Time `json:"expires"`
}

// NewAuth - Creates a new Auth
func NewAuth() Auth {
	authBytes := sha256.Sum256([]byte(strconv.Itoa(rand.Int())))
	authkey := base64.StdEncoding.EncodeToString([]byte(authBytes[:]))
	return Auth{
		Key:     authkey,
		Expires: time.Now().Add(24 * time.Hour),
	}
}
