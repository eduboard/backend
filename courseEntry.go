package eduboard

import (
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"time"
)

type CourseEntry struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id"`
	CourseID  bson.ObjectId `json:"courseID" bson:"courseID"`
	Date      time.Time     `json:"date" bson:"date"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	Message   string        `json:"message" bson:"message"`
	Pictures  []url.URL     `json:"pictures" bson:"pictures"`
	Published bool          `json:"published" bson:"published"`
}

type CourseEntryRepository interface {
	CourseEntryInserter
	CourseEntryOneFinder
	CourseEntryManyFinder
	CourseEntryUpdater
	CourseEntryDeleter
}

type CourseEntryInserter interface {
	Insert(course CourseEntry) error
}

type CourseEntryOneFinder interface {
	FindOneByID(id string) (error, CourseEntry)
}

type CourseEntryManyFinder interface {
	FindMany(query bson.M) (error, []CourseEntry)
}

type CourseEntryUpdater interface {
	Update(id string, update bson.M) error
}

type CourseEntryDeleter interface {
	Delete(id string) error
}

type CourseEntryService interface {
	StoreCourseEntry(entry *CourseEntry, cfu CourseFindUpdater) (err error, courseEntry *CourseEntry)
	StoreCourseEntryFiles(files [][]byte, id string, date time.Time) (err error, paths []string)
	UpdateCourseEntry(*CourseEntry) (*CourseEntry, error)
	DeleteCourseEntry(entryID string, courseID string, updater CourseUpdater) error
}
