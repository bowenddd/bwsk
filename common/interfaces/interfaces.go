package interfaces

import "seckill/common/entity"

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
}

type OrderServ interface {
	AddOrder(order *entity.Order) error

	GetOrderById(id uint) (entity.Order, error)

	GetOrdersByUID(uid uint) ([]entity.Order, error)

	GetOrdersByPID(pid uint) ([]entity.Order, error)

	DeleteOrder(id uint) error

	GetOrders() ([]entity.Order, error)
}
