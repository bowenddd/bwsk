package service

import (
	"fmt"
	"seckill/common/entity"
	"seckill/common/interfaces"
	dbrpc "seckill/dbservice/rpc"
	"seckill/registercenter/registerservice"
	"sync"
)

type UserServImpl struct {
	registerCenter *registerservice.RegisterCenter
}

func (p *UserServImpl) DbRpcCli() *dbrpc.UserRpcServCli {
	sc, err := p.registerCenter.GetDbClient()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	usc := sc.GetUserRpcServCli()
	return &usc
}

func (u *UserServImpl) AddUser(user *entity.User) error {
	return u.DbRpcCli().AddUser(user)
}

func (u *UserServImpl) GetUser(name string) (entity.User, error) {
	cli := u.DbRpcCli()
	return cli.GetUser(name)
}

func (u *UserServImpl) DeleteUser(name string) error {
	return u.DbRpcCli().DeleteUser(name)
}

func (u *UserServImpl) GetUsers() ([]entity.User, error) {
	return u.DbRpcCli().GetUsers()
}

var _ interfaces.UserServ = (*UserServImpl)(nil)

var userServ *UserServImpl

var userServOnce = new(sync.Once)

func GetUserService() interfaces.UserServ {

	userServOnce.Do(func() {
		center := registerservice.GetRegisterCenter()
		go func() {
			ch := make(chan error, 0)
			registerservice.GetRegisterCenter().Discovery(ch)
			<-ch
		}()
		userServ = &UserServImpl{
			registerCenter: center,
		}
	})
	return userServ
}
