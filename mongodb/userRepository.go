package mongodb

import (
	"github.com/eduboard/backend"
	"gopkg.in/mgo.v2"
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
	return nil
}

func (u *UserRepository) Find(id string) (error, *eduboard.User) {
	return nil, &eduboard.User{}
}
