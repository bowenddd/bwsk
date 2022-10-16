package service

import (
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/dbservice/rpc"
	"sync"
)

type UserServImpl struct {
	cli *rpc.UserRpcServCli
}

func (u *UserServImpl) AddUser(user *entity.User) error {
	return u.cli.AddUser(user)
}

func (u *UserServImpl) GetUser(name string) (entity.User, error) {
	return u.cli.GetUser(name)
}

func (u *UserServImpl) DeleteUser(name string) error {
	return u.cli.DeleteUser(name)
}

func (u *UserServImpl) GetUsers() ([]entity.User, error) {
	return u.cli.GetUsers()
}

var _ interfaces.UserServ = (*UserServImpl)(nil)

var userServ *UserServImpl

var userServOnce = new(sync.Once)

func GetUserService() interfaces.UserServ {
	cli, err := rpc.GetDbServRpcCli()
	if err != nil {
		return (*UserServImpl)(nil)
	}
	userServOnce.Do(func() {
		userRpcServCli := cli.GetUserRpcServCli()
		userServ = &UserServImpl{
			cli: &userRpcServCli,
		}
	})
	return userServ
}
