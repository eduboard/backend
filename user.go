package edubord

type User struct {
	Id       UserId `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Courses  []*CourseId
}

type UserId string

type UserRepository interface {
	Store(user *User) error
	Find(id UserId) (error, *User)
}

type UserService interface {
	CreateUser(u *User) (error, *User)
	Login(username string, password string) (error, *User)
}
