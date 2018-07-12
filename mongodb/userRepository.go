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
	if user.ID == "" {
		user.ID = bson.NewObjectId()
	}
	return u.c.Insert(user)
}

func (u *UserRepository) Find(id string) (error, eduboard.User) {
	result := eduboard.User{}

	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id"), result
	}
	if err := u.c.FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		return err, eduboard.User{}
	}

	return nil, result
}

func (u *UserRepository) FindMany(query bson.M) ([]eduboard.User, error) {
	result := []eduboard.User{}

	if err := u.c.Find(query).All(&result); err != nil {
		return []eduboard.User{}, err
	}

	return result, nil
}

func (u *UserRepository) FindBySessionID(sessionID string) (error, eduboard.User) {
	if sessionID == "" {
		return errors.New("not found"), eduboard.User{}
	}
	return u.findBy("sessionID", sessionID)
}

func (u *UserRepository) FindByEmail(email string) (error, eduboard.User) {
	if email == "" {
		return errors.New("not found"), eduboard.User{}
	}
	return u.findBy("email", email)
}

func (u *UserRepository) findBy(key string, value string) (error, eduboard.User) {
	result := eduboard.User{}

	if err := u.c.Find(bson.M{key: value}).One(&result); err != nil {
		return err, eduboard.User{}
	}

	return nil, result
}

func (u *UserRepository) IsIDValid(id string) bool {
	return bson.IsObjectIdHex(id)
}

func (u *UserRepository) FindMembers(members []string) (error, []eduboard.User) {
	result := []eduboard.User{}
	membersId := []bson.ObjectId{}

	for i := 0; i < len(members); i++ {
		membersId = append(membersId, bson.ObjectIdHex(members[i]))
	}

	if err := u.c.Find(bson.M{"_id": bson.M{"$in": membersId}}).All(&result); err != nil {
		return errors.New("error while fetching users"), []eduboard.User{}
	}

	return nil, result
}

func (u *UserRepository) UpdateSessionID(user eduboard.User) (error, eduboard.User) {
	return u.updateValue(user.ID.Hex(), bson.M{"$set": bson.M{"sessionID": user.SessionID}})
}

func (u *UserRepository) updateValue(id string, change bson.M) (error, eduboard.User) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id"), eduboard.User{}
	}

	err := u.c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, change)
	if err != nil {
		return err, eduboard.User{}
	}

	err, usr := u.Find(id)

	return err, usr
}
