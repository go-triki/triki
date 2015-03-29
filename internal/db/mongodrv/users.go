package mongo

import (
	"errors"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/models/user"
)

// usersC returns the users collection.
func usersC() *mgo.Collection {
	return adminSession.Copy().DB("").C(usersCName)
}

// usersSetup ensures that the users collection is setup correctly.
func usersSetup() {
	c := usersC()
	index := mgo.Index{
		Key:        []string{"usr"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		log.Fatalf("MongoDB ensureIndex `usr` on users failed: %s.\n", err.Error())
	}
}

// UserFind finds user with a given login/email.
func UserFind(usr string) (User, error) {
	var user user.T
	err := usersC().Find(bson.M{"usr": usr}).One(&user)
	return user, err
}

// UserFindByID finds user with a given _id.
func UserFindByID(id bson.ObjectId) (User, error) {
	var user user.T
	err := usersC().Find(bson.M{"_id": id}).One(&user)
	return user, err
}

// UserSignup signs given user up, pending email verification.
// FIXME: write email verification
func UserSignup(login, pass string) error {
	_, err := UserFindByLogin(login)
	if err == nil { // FIXME: is this check enough?
		return errors.New("User already exists")
	}
	usr := User{
		Login:    login,
		Password: pass,
	}
	err = userNew(&usr)
	return err
}
