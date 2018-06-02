package eduboard

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name         string        `json:"name" bson:"name"`
	Surname      string        `json:"surname" bson:"surname"`
	Email        string        `json:"email" bson:"email"`
	PasswordHash string        `json:"password" bson:"password"`
	SessionID    string        `json:"sessionId" bson:"sessionId"`
	Courses      []string      `json:"courses" bson:"courses"`
}

type UserRepository interface {
	Store(user *User) error
	Find(id string) (error, *User)
	FindByEmail(email string) (error, *User)
	FindBySessionID(sessionID string) (error, *User)
	UpdateSessionID(user *User) (error, *User)
}

type UserService interface {
	CreateUser(u *User, password string) (error, *User)
	GetUser(id string) (error, *User)
	UserAuthenticationProvider
}

type UserAuthenticationProvider interface {
	Login(email string, password string) (error, *User)
	Logout(sessionID string) error
	CheckAuthentication(sessionID string) (err error, ok bool)
}
