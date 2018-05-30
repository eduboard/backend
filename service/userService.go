package service

import (
	"errors"
	"github.com/eduboard/backend"
)

type UserService struct {
	r eduboard.UserRepository
	a eduboard.UserAuthenticator
}

func NewUserService(repository eduboard.UserRepository, authenticator eduboard.UserAuthenticator) *UserService {
	return &UserService{
		r: repository,
		a: authenticator,
	}
}

func (uS *UserService) CreateUser(user *eduboard.User) (error, *eduboard.User) {
	
	hashedPassword := uS.a.HashAndSalt(user.PasswordHash)
	user.PasswordHash = hashedPassword

	err := uS.r.Store(user)
	
	return err, user
}

func (uS *UserService) GetUser(id string) (err error, user *eduboard.User) {
	return uS.r.Find(id)
}

func (uS *UserService) Login(username string, password string) (error, *eduboard.User) {
	err, user := uS.r.FindByUsername(username)
	if err != nil {
		return err, &eduboard.User{}	
	}

	if !uS.a.CompareWithHash(password, user.PasswordHash) {
		return errors.New("invalid password"), &eduboard.User{}	
	}

	user.AccessToken = uS.a.CreateAccessToken()
	err, user = uS.r.UpdateAccessToken(user)
	if err != nil {
		return err, &eduboard.User{}
	}

	return nil, user
}

func (uS *UserService) Logout(accessToken string) (error) {
	err, user := uS.r.FindByAccessToken(accessToken)
	if err != nil {
		return err
	}

	user.AccessToken = ""
	uS.r.UpdateAccessToken(user)

	return nil
}
