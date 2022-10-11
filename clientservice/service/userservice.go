package service

import (
	"fmt"
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/dbservice/rpc"
	"seckill/seetings"
	"sync"
	"time"
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

var dbServRpcCliOnce = new(sync.Once)

var rpcCli *rpc.ServClient

func GetDbServRpcCli() (*rpc.ServClient, error) {
	dbServRpcCliOnce.Do(func() {
		addr, timeout, err := getRpcSettings()
		if err != nil {
			return
		}
		cli, err := rpc.NewServClient(addr, time.Duration(timeout)*time.Second)
		if err != nil {
			return
		}
		rpcCli = &cli
	})
	if rpcCli == nil {
		return rpcCli, fmt.Errorf("get db service rpc client error")
	}
	return rpcCli, nil
}

func GetUserService() interfaces.UserServ {
	cli, err := GetDbServRpcCli()
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

func getRpcSettings() (addr string, timeout int, err error) {
	setting, err := seetings.GetSetting()
	if err != nil {
		fmt.Println("get setting error in clientservice!")
		return
	}
	port := setting.RPC.DbServPort
	addr = fmt.Sprintf("localhost%s", port)
	timeout = setting.RPC.Timeout
	return
}
