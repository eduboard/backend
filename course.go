package edubord

type Course struct {
	Id          CourseId `json:"id" bson:"_id"`
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	Members     []UserId `json:"members" bson:"members"`
}

type CourseId string

type CourseRepository interface {
	Store(course *Course) error
	Find(id CourseId) (error, *Course)
	FindAll() (error, []*Course)
}

type CourseService interface {
	GetAllCourses() (err error, courses []*Course)
	GetCourse(id CourseId) (err error, course *Course)
}
