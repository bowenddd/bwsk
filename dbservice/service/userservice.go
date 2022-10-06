package service

import (
	"seckill/dbservice/store"
	"seckill/entity"
	"sync"
)

type UserServ interface {
	AddUser(user *entity.User) error

	GetUser(name string) (entity.User, error)

	DeleteUser(name string) error

	GetUsers() ([]entity.User, error)
}

type UserServImpl struct {
	store store.UserStore
}

var userServOnce sync.Once

var userServ UserServ

func init() {
	userServOnce = sync.Once{}
}

func GetUserServ() UserServ {
	userServOnce.Do(func() {
		userServ = &UserServImpl{
			store: store.Mysql.NewUserStore(),
		}
	})
	return userServ
}

func (u UserServImpl) AddUser(user *entity.User) error {
	return u.store.Create(user)
}

func (u UserServImpl) GetUser(name string) (entity.User, error) {
	return u.store.FindByName(name)
}

func (u UserServImpl) DeleteUser(name string) error {
	return u.store.DeleteByName(name)
}

func (u UserServImpl) GetUsers() ([]entity.User, error) {
	return u.store.List()
}

var _ UserServ = (*UserServImpl)(nil)
