package service

import (
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/dbservice/store"
	"strings"
)

type PermServImpl struct {
	store store.PermStore
}

func (p *PermServImpl) AddRole(role *entity.Role) error {
	return p.store.AddRole(role)
}

func (p *PermServImpl) AddPerm(perm *entity.Perm) error {
	return p.store.AddPerm(perm)
}

func (p *PermServImpl) RoleList() ([]entity.Role, error) {
	return p.store.RoleList()
}

func (p *PermServImpl) PermList() ([]entity.Perm, error) {
	return p.store.PermList()
}

func (p *PermServImpl) SetRole(userId, roleId int) error {
	return p.store.SetRole(userId, roleId)
}

func (p *PermServImpl) GetPerm(uid uint) (string, error) {
	perms, err := p.store.GetPerm(uid)
	if err != nil {
		return "", err
	}
	return strings.Join(perms, ","), nil
}

var _ interfaces.PermServ = (*PermServImpl)(nil)

func GetPermServ() interfaces.PermServ {
	permStore := store.Mysql.NewPermStore()
	permServImpl := &PermServImpl{
		store: permStore,
	}
	return permServImpl
}
