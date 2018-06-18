package service

import (
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/auth"
	"github.com/pkg/errors"
	"time"
)

type UserService struct {
	r eduboard.UserRepository
	a Authenticator
}

type Authenticator interface {
	Hash(password string) (string, error)
	CompareHash(hashedPassword string, plainPassword string) (bool, error)
	SessionID() string
}

func NewUserService(repository eduboard.UserRepository) *UserService {
	return &UserService{
		r: repository,
		a: &auth.Authenticator{},
	}
}

func (uS *UserService) CreateUser(user *eduboard.User, password string) (error, *eduboard.User) {
	err, _ := uS.r.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already exists"), &eduboard.User{}
	}

	if len(password) < 8 {
		return errors.New("password too short"), &eduboard.User{}
	}

	hashedPassword, err := uS.a.Hash(password)
	if err != nil {
		return errors.Wrap(err, "error hashing password"), &eduboard.User{}
	}

	user.PasswordHash = hashedPassword
	user.SessionID = uS.a.SessionID()
	user.SessionExpires = time.Time{}.Add(24 * time.Hour)

	err = uS.r.Store(user)
	if err != nil {
		return errors.Wrap(err, "error storing user"), &eduboard.User{}
	}
	return nil, user
}

func (uS *UserService) GetUser(id string) (err error, user *eduboard.User) {
	return uS.r.Find(id)
}

func (uS *UserService) Login(email string, password string) (error, *eduboard.User) {
	err, user := uS.r.FindByEmail(email)
	if err != nil {
		return errors.Wrap(err, "error finding user by email"), &eduboard.User{}
	}

	ok, err := uS.a.CompareHash(user.PasswordHash, password)
	if err != nil {
		return errors.Wrap(err, "error comparing hash"), &eduboard.User{}
	}
	if !ok {
		return errors.New("invalid password"), &eduboard.User{}
	}

	user.SessionID = uS.a.SessionID()
	err, user = uS.r.UpdateSessionID(user)
	if err != nil {
		return errors.Wrap(err, "error updating sessionID"), &eduboard.User{}
	}

	return nil, user
}

func (uS *UserService) Logout(sessionID string) error {
	err, user := uS.r.FindBySessionID(sessionID)
	if err != nil {
		return errors.Wrap(err, "error finding by sessionID")
	}

	user.SessionID = ""
	uS.r.UpdateSessionID(user)

	return nil
}

func (uS *UserService) CheckAuthentication(sessionID string) (err error, id string) {
	err, user := uS.r.FindBySessionID(sessionID)
	if err != nil {
		return errors.Wrap(err, "error finding by sessionID"), ""
	}
	if user.SessionExpires.After(time.Now()) {
		return errors.New("session expired"), ""
	}

	return nil, user.ID.Hex()
}
