package db

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

const (
	// minimal password length
	usersMinPassLen = 4
)

// usersC returns the users collection.
func usersC() *mgo.Collection {
	return adminSession.Copy().DB("").C(usersCName)
}

// User type stores user information (e.g. for authentication), also for MongoDB and JSON.
type User struct {
	ID       bson.ObjectId `json:"id"             bson:"_id"`
	Login    string        `json:"login"          bson:"login"`
	Password string        `json:"pass,omitempty" bson:"-"`
	PassHash []byte        `json:"-"              bson:"pass"`
	Name     string        `json:"name"           bson:"name"`
}

// usersSetup ensures that the users collection is setup correctly.
func usersSetup() {
	c := usersC()
	index := mgo.Index{
		Key:        []string{"login"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		log.Fatalf("MongoDB ensureIndex login on users failed: %s.\n", err)
	}
}

// UserFindByLogin finds user with a given login.
func UserFindByLogin(login string) (User, error) {
	var user User
	err := usersC().Find(bson.M{"login": login}).One(&user)
	return user, err
}

// UserFindByID finds user with a given login.
func UserFindByID(id bson.ObjectId) (User, error) {
	var user User
	err := usersC().Find(bson.M{"_id": id}).One(&user)
	return user, err
}

// UserAuthenticate checks if the given credentials can be authenticated.
// Returns (user, token, error).
func UserAuthenticate(login, pass string) (*User, string, error) {
	usr, err := UserFindByLogin(login)
	if err != nil {
		return nil, "", err
	}
	err = bcrypt.CompareHashAndPassword(usr.PassHash, []byte(usersPassPrefix+pass))
	if err != nil {
		return nil, "", err
	}
	token, err := tokenNew(usr.ID)
	if err != nil {
		return nil, "", err
	}
	return &usr, token, err
}

// userCheck carries out some sanity checks on the given user (pass length, is login an email address).
func userCheck(usr *User) error {
	// password length
	if len(usr.Password) < usersMinPassLen {
		return errors.New("Password too short")
	}
	// does login look like an email?
	if len(usr.Login) < 1 {
		return errors.New("Login cannot be empty")
	}
	at := strings.Index(usr.Login, "@")
	if at < 1 || at == len(usr.Login)-1 {
		return errors.New("Login needs to be an email address")
	}
	if len(usr.Name) < 1 {
		return errors.New("Name cannot be empty")
	}
	// checks passed
	return nil
}

// UserNew adds new user to the database. Returns nil on success.
// User.Password is hashed into User.PassHash.
// ID is generated.
func userNew(usr *User) error {
	if usr.Name == "" {
		usr.Name = usr.Login
	}
	err := userCheck(usr)
	if err != nil {
		return err
	}
	usr.ID = bson.NewObjectId()
	usr.PassHash, err = bcrypt.GenerateFromPassword([]byte(usersPassPrefix+usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = usersC().Insert(usr)
	return err
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
