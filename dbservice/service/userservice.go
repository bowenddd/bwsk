package service

import (
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/dbservice/store"
	"sync"
)

type UserServImpl struct {
	store store.UserStore
}

var userServOnce sync.Once

var userServ interfaces.UserServ

func init() {
	userServOnce = sync.Once{}
}

func GetUserServ() interfaces.UserServ {
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

var _ interfaces.UserServ = (*UserServImpl)(nil)
