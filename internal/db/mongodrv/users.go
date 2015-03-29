package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/user"
)

// usersC returns the users collection.
func usersC() *mgo.Collection {
	return adminSession.Copy().DB("").C(usersCName)
}

// UserInsert inserts the usr into the DB.
func UserInsert(usr *user.T) *log.Error {
	err := usersC().Insert(usr)
	if err != nil {
		return log.InternalServerErr(err)
	}
	return nil
}

// UserFind finds user with a given login/email.
func UserFind(login string) (*user.T, *log.Error) {
	var usr user.T
	err := usersC().Find(bson.M{"usr": login}).One(&usr)
	if err != nil {
		return nil, log.DBNotFoundErr(err)
	}
	return &usr, nil
}

// UserFindByID finds user with a given _id.
func UserFindByID(id bson.ObjectId) (*user.T, *log.Error) {
	var usr user.T
	err := usersC().Find(bson.M{"_id": id}).One(&usr)
	if err != nil {
		return nil, log.DBNotFoundErr(err)
	}
	return &usr, nil
}

// UserExists checks if user with login/email == usr exists in the DB.
func UserExists(login string) (bool, *log.Error) {
	n, err := usersC().Find(bson.M{"usr": login}).Count()
	if err != nil {
		return false, log.InternalServerErr(err)
	}
	return n >= 1, nil
}
