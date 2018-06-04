package eduboard

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name         string        `json:"name" bson:"name"`
	Surname      string        `json:"surname" bson:"surname"`
	Email        string        `json:"email" bson:"email"`
	PasswordHash string        `json:"password" bson:"password"`
	SessionId    string        `json:"sessionId" bson:"sessionId"`
	Courses      []string      `json:"courses" bson:"courses"`
}

type UserRepository interface {
	Store(user *User) error
	Find(id string) (error, *User)
	FindByEmail(email string) (error, *User)
	FindBySessionId(sessionId string) (error, *User)
	UpdateSessionId(user *User) (error, *User)
}

type UserService interface {
	CreateUser(u *User, password string) (error, *User)
	GetUser(id string) (error, *User)
	UserAuthenticationProvider
}

type UserAuthenticationProvider interface {
	Login(email string, password string) (error, *User)
	Logout(sessionId string) error
	CheckAuthentication(sessionId string) (err error, ok bool)
}
