package service

import (
	"errors"
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/auth"
)

type UserService struct {
	r eduboard.UserRepository
	a Authenticator
}

type Authenticator interface {
	Hash(password string) (string, error)
	CompareHash(hashedPassword string, plainPassword string) (bool, error)
	SessionId() string
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
		return err, &eduboard.User{}
	}
	user.PasswordHash = hashedPassword
	user.SessionId = uS.a.SessionId()

	err = uS.r.Store(user)
	if err != nil {
		return err, &eduboard.User{}
	}
	return nil, user
}

func (uS *UserService) GetUser(id string) (err error, user *eduboard.User) {
	return uS.r.Find(id)
}

func (uS *UserService) Login(email string, password string) (error, *eduboard.User) {
	err, user := uS.r.FindByEmail(email)
	if err != nil {
		return err, &eduboard.User{}
	}

	ok, err := uS.a.CompareHash(user.PasswordHash, password)
	if err != nil {
		return err, &eduboard.User{}
	}
	if !ok {
		return errors.New("invalid password"), &eduboard.User{}
	}

	user.SessionId = uS.a.SessionId()
	err, user = uS.r.UpdateSessionId(user)
	if err != nil {
		return err, &eduboard.User{}
	}

	return nil, user
}

func (uS *UserService) Logout(sessionId string) error {
	err, user := uS.r.FindBySessionId(sessionId)
	if err != nil {
		return err
	}

	user.SessionId = ""
	uS.r.UpdateSessionId(user)

	return nil
}

func (uS *UserService) CheckAuthentication(sessionId string) (err error, ok bool) {
	err, _ = uS.r.FindBySessionId(sessionId)
	if err != nil {
		return err, false
	}
	return nil, true
}
