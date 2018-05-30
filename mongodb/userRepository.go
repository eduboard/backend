package mongodb

import (
	"errors"
	"github.com/eduboard/backend"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
	c *mgo.Collection
}

func newUserRepository(database *mgo.Database) *UserRepository {
	collection := database.C("user")
	return &UserRepository{
		c: collection,
	}
}

func (u *UserRepository) Store(user *eduboard.User) error {
	return u.c.Insert(user)
}

func (u *UserRepository) Find(id string) (error, *eduboard.User) {
	result := eduboard.User{}

	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id"), &result
	}
	if err := u.c.FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		return err, &eduboard.User{}
	}

	return nil, &result
}

func (u *UserRepository) FindByAccessToken(accessToken string) (error, *eduboard.User) {
	return u.findBy("accessToken", accessToken)
}

func (u *UserRepository) FindByUsername(username string) (error, *eduboard.User) {
	return u.findBy("username", username)
}

func (u *UserRepository) findBy(key string, value string) (error, *eduboard.User) {
	result := eduboard.User{}

	if err:= u.c.Find(bson.M{key: value}).One(&result); err != nil {
		return err, &eduboard.User{}
	}

	return nil, &result
}

func (u *UserRepository) UpdateAccessToken(user *eduboard.User) (error, *eduboard.User){
	return u.updateValue(user.Id, "accessToken", user.AccessToken)
}

func (u *UserRepository)updateValue(id bson.ObjectId, key string, value string) (error, *eduboard.User){
	user := bson.M{"_id": id}

	change := bson.M{"$set": bson.M{key: value}}
	err := u.c.Update(user, change)
	if err != nil {
		return err, &eduboard.User{}
	}

	err, usr := u.Find(id.Hex())

	return err, usr
}
