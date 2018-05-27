package service

import "github.com/eduboard/backend"

type UserService struct {
	r eduboard.UserRepository
}

func NewUserService(repository eduboard.UserRepository) *UserService {
	return &UserService{
		r: repository,
	}
}

func (uS *UserService) CreateUser(user *eduboard.User) (error, *eduboard.User) {
	return nil, &eduboard.User{}
}
func (uS *UserService) Login(username string, password string) (error, *eduboard.User) {
	return nil, &eduboard.User{}
}
