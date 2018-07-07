package prefService

import (
	"github.com/eduboard/backend"

)

func New(repository eduboard.PrefRepository) PrefService {
	return PrefService{
		PR: repository,
	}
}

type PrefService struct {
	PR eduboard.PrefRepository
}

func (pS PrefService) CreatePref(c *eduboard.Pref) (*eduboard.Pref, error) {
	return c, nil
}