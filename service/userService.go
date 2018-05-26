package service

import "github.com/eduboard/backend"

type UserService struct {
	r edubord.UserRepository
}

func NewUserService(repository edubord.UserRepository) *UserService {
	return &UserService{
		r: repository,
	}
}

func (u *UserService) CreateUser(user *edubord.User) (error, *edubord.User) {
	return nil, &edubord.User{}
}
func (u *UserService) Login(username string, password string) (error, *edubord.User) {
	return nil, &edubord.User{}
}
