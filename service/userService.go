package service

import (
	"errors"
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/auth"
)

type UserService struct {
	r eduboard.UserRepository
	a auth.Authenticator
}

type AuthenticationProvider interface {
	Hash(password string) (string, error)
	CompareHash(hashedPassword string, plainPassword string) (bool, error)
	SessionId() string
}

func NewUserService(repository eduboard.UserRepository) *UserService {
	return &UserService{
		r: repository,
		a: auth.Authenticator{},
	}
}

func (uS *UserService) CreateUser(user *eduboard.User) (error, *eduboard.User) {
	hashedPassword, err := uS.a.Hash(user.PasswordHash)
	if err != nil {
		return err, &eduboard.User{}
	}
	user.PasswordHash = hashedPassword

	err = uS.r.Store(user)
	if err != nil {
		return err, &eduboard.User{}
	}
	return nil, user
}

func (uS *UserService) GetUser(id string) (err error, user *eduboard.User) {
	return uS.r.Find(id)
}

func (uS *UserService) Login(username string, password string) (error, *eduboard.User) {
	err, user := uS.r.FindByUsername(username)
	if err != nil {
		return err, &eduboard.User{}
	}

	ok, err := uS.a.CompareHash(password, user.PasswordHash)
	if err != nil {
		return err, &eduboard.User{}
	}
	if !ok {
		return errors.New("invalid password"), &eduboard.User{}
	}

	user.AccessToken = uS.a.SessionId()
	err, user = uS.r.UpdateAccessToken(user)
	if err != nil {
		return err, &eduboard.User{}
	}

	return nil, user
}

func (uS *UserService) Logout(sessionId string) error {
	err, user := uS.r.FindByAccessToken(sessionId)
	if err != nil {
		return err
	}

	user.AccessToken = ""
	uS.r.UpdateAccessToken(user)

	return nil
}

func (uS *UserService) CheckAuthentication(sessionId string) (err error, ok bool) {
	err, _ = uS.r.FindByAccessToken(sessionId)
	if err != nil {
		return err, false
	}
	return nil, true
}
