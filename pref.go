package eduboard

import (
	"gopkg.in/mgo.v2/bson"
)

type Pref struct {
	ID bson.ObjectId `json:"id,omitempty" bson:"_id"`
	CourseID  string `json:"courseID" bson:"courseID"`
	UserID    string `json:"userID" bson:"userID"`
	PartnerID string `json:"partnerID" bson:"partnerID"`
	Value     int    `json:"value" bson:"value"`
}

type PrefInserter interface {
	Insert(course *Course) error
}

type PrefFinder interface {
	FindMany(query bson.M) ([]Pref, error)
}

type PrefRepository interface {
	PrefInserter
	PrefFinder
}

type PrefService interface {
	CreatePref(c *Pref) (*Pref, error)
	GetPrefByUserAndCourse(id string, cef CourseEntryManyFinder) (courses []Pref, err error)
}
