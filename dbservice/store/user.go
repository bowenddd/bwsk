package store

import (
	"gorm.io/gorm"
	"seckill/entity"
)

type UserStore struct {
	db *gorm.DB
}

func (us *UserStore) Create(user *entity.User) error {
	return us.db.Create(user).Error
}

func (us *UserStore) FindByName(name string) (*entity.User, error) {
	u := &entity.User{}
	res := us.db.Where("name = ?", name).First(u)
	return u, res.Error
}

func (us *UserStore) DeleteByName(name string) error {
	exec := us.db.Exec("delete from user where name = ?", name)
	return exec.Error
}
