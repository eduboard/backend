package eduboard

import (
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"time"
)

type User struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	Name           string        `json:"name" bson:"name"`
	Surname        string        `json:"surname" bson:"surname"`
	Email          string        `json:"email" bson:"email"`
	PasswordHash   string        `json:"password" bson:"password"`
	SessionID      string        `json:"sessionID" bson:"sessionID"`
	SessionExpires time.Time     `json:"sessionExpires" bson:"sessionExpires"`
	Courses        []string      `json:"courses" bson:"courses"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	Picture        url.URL       `json:"profilePicture" bson:"profilePicture"`
}

type UserFinder interface {
	FindMembers(members []string) (error, []User)
}

type UserRepository interface {
	Store(user *User) error
	Find(id string) (error, User)
	FindByEmail(email string) (error, User)
	FindBySessionID(sessionID string) (error, User)
	IsIDValid(id string) bool
	UpdateSessionID(user User) (error, User)
	UserFinder
}

type UserService interface {
	CreateUser(u *User, password string) (error, User)
	GetUser(id string) (error, User)
	GetMyCourses(id string, cS CourseManyFinder, cEMF CourseEntryManyFinder) (error, []Course)
	UserAuthenticationProvider
}

type UserAuthenticationProvider interface {
	Login(email string, password string) (error, User)
	Logout(sessionID string) error
	CheckAuthentication(sessionID string) (err error, userID string)
}
