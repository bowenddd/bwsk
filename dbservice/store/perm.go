package store

import (
	"errors"
	"seckill/common/entity"

	"gorm.io/gorm"
)

type PermStore interface {

	RoleList() ([]entity.Role, error)

	PermList() ([]entity.Perm, error)

	AddRole(*entity.Role) error

	AddPerm(*entity.Perm) error

	GetPerm(uid uint) ([]string, error)

	SetRole(userId,roleId int) error
}

type PermOp struct {
	db *gorm.DB
}

func (p *PermOp) DB() *gorm.DB {
	return p.db
}

func (p *PermOp) RoleList() ([]entity.Role, error) {
	var roles []entity.Role
	res := p.DB().Find(&roles)
	return roles, res.Error
}

func (p *PermOp) PermList() ([]entity.Perm, error) {
	var perms []entity.Perm
	res := p.DB().Find(&perms)
	return perms, res.Error
}

func (p *PermOp) SetRole(userId,roleId int) error {
	roleIds := []int{}
	p.DB().Table("user_role").Where("user_id = ?", userId).Select("role_id").Scan(&roleIds)
	for _, id := range roleIds {
		if id == roleId {
			return errors.New("user already has this role")
		}
	}
	res := p.DB().Exec("insert into user_role (user_id,role_id) values (?,?)", userId, roleId)
	return res.Error
}

func (p *PermOp) AddRole(role *entity.Role) error {
	return p.DB().Create(role).Error
}

func (p *PermOp) AddPerm(perm *entity.Perm) error {
	return p.DB().Create(perm).Error
}

func (p *PermOp) GetPerm(uid uint) ([]string, error) {
	var perms []string
	res := p.DB().Table("perm").Select("path").Joins("join role_perm on perm.id = role_perm.perm_id").Joins("join user_role on role_perm.role_id = user_role.role_id").Where("user_role.user_id = ?", uid).Scan(&perms)
	return perms, res.Error
}
