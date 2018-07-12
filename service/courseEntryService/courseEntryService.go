package courseEntryService

import (
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/upload"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

func New(repository eduboard.CourseEntryRepository) CourseEntryService {
	return CourseEntryService{
		ER: repository,
		u:  &upload.Uploader{},
	}
}

type Uploader interface {
	UploadFile(file []byte, course string, filename string) (string, error)
}

type CourseEntryService struct {
	ER eduboard.CourseEntryRepository
	u  Uploader
}

func (cES CourseEntryService) StoreCourseEntry(entry *eduboard.CourseEntry, cfu eduboard.CourseFindUpdater) (error, *eduboard.CourseEntry) {
	courseID := entry.CourseID.Hex()
	err, course := cfu.FindOneByID(courseID)
	if err != nil {
		return errors.Wrapf(err, "error finding course with ID %d", courseID), &eduboard.CourseEntry{}
	}

	entryID := bson.NewObjectId()
	entry.ID = entryID
	err = cES.ER.Insert(*entry)
	if err != nil {
		return errors.Wrap(err, "error inserting courseEntry"), &eduboard.CourseEntry{}
	}

	courseID = course.ID.Hex()
	err, _ = cfu.Update(courseID, bson.M{"$push": bson.M{"entryIDs": entryID}})
	if err != nil {
		return errors.Wrapf(err, "error updating course with ID %s", courseID), &eduboard.CourseEntry{}
	}

	return nil, entry
}

func (cES CourseEntryService) StoreCourseEntryFiles(files [][]byte, id string, date time.Time) ([]string, error) {
	paths := []string{}
	for key, file := range files {
		url, err := cES.u.UploadFile(file, id, date.String()+"_"+strconv.Itoa(key))
		if err != nil {
			return []string{}, err
		}
		paths = append(paths, url)
	}
	return paths, nil
}

func (cES CourseEntryService) UpdateCourseEntry(*eduboard.CourseEntry) (*eduboard.CourseEntry, error) {
	return &eduboard.CourseEntry{}, nil
}

func (cES CourseEntryService) DeleteCourseEntry(entryID string, courseID string, cfu eduboard.CourseUpdater) error {
	err, entry := cES.ER.FindOneByID(entryID)
	if err != nil {
		return errors.Errorf("could not find entry with id %s", entryID)
	}

	if entry.CourseID.Hex() != courseID {
		return errors.Errorf("entry with ID %s does not belong to course with ID %s", entryID, courseID)
	}

	if err := cES.ER.Delete(entryID); err != nil {
		return errors.Wrapf(err, "error deleting courseEntry with ID %s", entryID)
	}

	if err, _ = cfu.Update(courseID, bson.M{"$pull": bson.M{"entryIDs": entryID}}); err != nil {
		return err
	}

	return nil
}
