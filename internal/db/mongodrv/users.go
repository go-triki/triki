package mongo

import (
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/user"
)

// usersC returns the users collection.
func usersC(cx context.Context) (*mgo.Collection, *log.Error) {
	sess, err := getSession(cx, sessKey)
	if err != nil {
		return nil, err
	}
	return sess.DB("").C(usersCName), nil
}

// UserInsert inserts the usr into the DB.
func UserInsert(cx context.Context, usr *user.T) *log.Error {
	col, er := usersC(cx)
	if er != nil {
		return er
	}
	err := col.Insert(usr)
	if err != nil {
		return log.InternalServerErr(err)
	}
	return nil
}

// UserFind finds user with a given login/email.
func UserFind(cx context.Context, login string) (*user.T, *log.Error) {
	var usr user.T
	col, er := usersC(cx)
	if er != nil {
		return nil, er
	}
	err := col.Find(bson.M{"usr": login}).One(&usr)
	if err != nil {
		return nil, log.DBNotFoundErr(err)
	}
	return &usr, nil
}

// UserFindByID finds user with a given _id.
func UserFindByID(cx context.Context, id bson.ObjectId) (*user.T, *log.Error) {
	var usr user.T
	col, er := usersC(cx)
	if er != nil {
		return nil, er
	}
	err := col.Find(bson.M{"_id": id}).One(&usr)
	if err != nil {
		return nil, log.DBNotFoundErr(err)
	}
	return &usr, nil
}

// UserExists checks if user with login/email == usr exists in the DB.
func UserExists(cx context.Context, login string) (bool, *log.Error) {
	col, er := usersC(cx)
	if er != nil {
		return false, er
	}
	n, err := col.Find(bson.M{"usr": login}).Count()
	if err != nil {
		return false, log.InternalServerErr(err)
	}
	return n >= 1, nil
}
