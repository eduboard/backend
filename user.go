package eduboard

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username     string        `json:"username" bson:"username"`
	Name         string        `json:"name" bson:"name"`
	Surname      string        `json:"surname" bson:"surname"`
	Email        string        `json:"email" bson:"email"`
	PasswordHash string        `json:"password" bson:"password"`
	AccessToken  string        `json:"accessToken" bson:"accessToken"`
	Courses      []string      `json:"courses" bson:"courses"`
}

type UserRepository interface {
	Store(user *User) error
	Find(id string) (error, *User)
	FindByUsername(username string) (error, *User)
	FindByAccessToken(accessToken string) (error, *User)
	UpdateAccessToken(user *User) (error, *User)
}

type UserService interface {
	CreateUser(u *User) (error, *User)
	GetUser(id string) (error, *User)
	UserAuthenticationProvider
}

type UserAuthenticationProvider interface {
	Login(username string, password string) (error, *User)
	Logout(sessionId string) error
	CheckAuthentication(sessionId string) (err error, ok bool)
}
