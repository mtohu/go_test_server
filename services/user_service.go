package services

import (
	"gotest/datamodels"
	"gotest/repositories"
)


type IUserService interface {
	//GetAll() []datamodels.User
	//GetByID(id int64) (datamodels.User, bool)
	GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool)
	//DeleteByID(id int64) bool

	//Update(id int64, user datamodels.User) (datamodels.User, error)
	//UpdatePassword(id int64, newPassword string) (datamodels.User, error)
	//UpdateUsername(id int64, newUsername string) (datamodels.User, error)

	//Create(userPassword string, user datamodels.User) (datamodels.User, error)
}

type UserService struct {
	repo repositories.IUserRepository
}

func NewUserService() IUserService {
	return &UserService{repo: repositories.NewUserRepository()}
}

func (s *UserService) GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool) {
	if username == "" || userPassword == "" {
		return datamodels.User{}, false
	}

	return s.repo.Select(func(m datamodels.User) bool {
		if m.Username == username {
			return true
		}
		return false
	})
}
