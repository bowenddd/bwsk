package interfaces

import (
	"seckill/common/entity"
	"time"
)

type UserServ interface {
	AddUser(user *entity.User) error

	GetUser(name string) (entity.User, error)

	DeleteUser(name string) error

	GetUsers() ([]entity.User, error)
}

type ProductServ interface {
	AddProduct(product *entity.Product) error

	GetProduct(name string) (entity.Product, error)

	GetProducts() ([]entity.Product, error)

	DeleteProduct(name string) error

	SetStock(id uint, num int) error

	GetStock(int uint) (int, error)
}

type OrderServ interface {
	AddOrder(order *entity.Order, method string) error

	GetOrderById(id uint) (entity.Order, error)

	GetOrdersByUID(uid uint) ([]entity.Order, error)

	GetOrdersByPID(pid uint) ([]entity.Order, error)

	DeleteOrder(id uint) error

	GetOrders() ([]entity.Order, error)

	ClearOrders() error
}

type CacheServ interface {
	SetStock(id uint, num int, exp time.Duration) error

	GetStock(id uint) (int, error)

	CreateOrder(order *entity.Order, method string) error

	Lock(key string, ex time.Duration) (bool, error)

	UnLock (key string) int64
}
