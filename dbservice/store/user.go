package store

import (
	"gorm.io/gorm"
	"seckill/entity"
)

type UserStore interface {
	Create(user *entity.User) error

	FindByName(name string) (entity.User, error)

	DeleteByName(name string) error

	List() ([]entity.User, error)
}

var _ UserStore = (*UserOp)(nil)

type UserOp struct {
	db *gorm.DB
}

func (us *UserOp) DB() *gorm.DB {
	return us.db
}

func (us *UserOp) Create(user *entity.User) error {
	return us.DB().Create(user).Error
}

func (us *UserOp) FindByName(name string) (entity.User, error) {
	u := entity.User{}
	res := us.DB().Where("name = ?", name).First(&u)
	return u, res.Error
}

func (us *UserOp) DeleteByName(name string) error {
	exec := us.DB().Exec("delete from user where name = ?", name)
	return exec.Error
}

func (us *UserOp) List() ([]entity.User, error) {
	var users []entity.User
	res := us.DB().Find(&users)
	return users, res.Error
}
