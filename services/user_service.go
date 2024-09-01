package services

import (
	"golang.org/x/crypto/bcrypt"
	"imooc-product/datamodels"
	"imooc-product/repositories"
)

type IUserService interface {
	IsPwdSuccess(username string, password string) (user *datamodels.User, isOK bool)
	AddUser(user *datamodels.User) (userId int64, err error)
	//
}
type UserService struct {
	UserRepository repositories.IUserRepository
}

func ValidatePassword(userPassword string, hasPwd string) (bool, error) {
	if userPassword == hasPwd {
		return true, nil
	} else {
		return false, nil
	}
}
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}
func (u *UserService) IsPwdSuccess(username string, password string) (user *datamodels.User, isOK bool) {
	user, err := u.UserRepository.Select(username)
	if err != nil {
		return nil, false
	}
	isOK, err = ValidatePassword(password, user.HashPassword)
	if err != nil {
		return nil, false
	}
	return user, isOK
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	password, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return 0, err
	}
	user.HashPassword = string(password)
	id, err := u.UserRepository.Insert(user)
	return id, err
}

func NewService(userRepository repositories.IUserRepository) IUserService {
	return &UserService{UserRepository: userRepository}
}
