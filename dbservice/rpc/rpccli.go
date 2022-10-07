package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"seckill/dbservice/service"
	"seckill/entity"
	pb "seckill/rpc/dbservice"
	"time"
)

type ServClient struct {
	conn    *grpc.ClientConn
	timeout time.Duration
}

type UserRpcServCli struct {
	rpcClient pb.UserServClient
	timeout   time.Duration
}

type OrderRpcServCli struct {
	rpcClient pb.OrderServClient
	timeout   time.Duration
}

type ProductRpcServCli struct {
	rpcClient pb.ProductServClient
	timeout   time.Duration
}

func NewServClient(addr string, timeout time.Duration) (ServClient, error) {
	var conn *grpc.ClientConn
	var err error
	if timeout == -1 {
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		timeout = math.MaxInt64
	} else {
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(timeout))
	}
	if err != nil {
		fmt.Println("grpc dial error! in new serv client")
		return ServClient{}, err
	}
	return ServClient{
		conn:    conn,
		timeout: timeout,
	}, nil
}

func (s *ServClient) GetUserRpcServCli() UserRpcServCli {
	client := pb.NewUserServClient(s.conn)
	userRpcCli := UserRpcServCli{
		rpcClient: client,
		timeout:   s.timeout,
	}
	return userRpcCli
}

func (s *ServClient) GetOrderRpcServCli() OrderRpcServCli {
	client := pb.NewOrderServClient(s.conn)
	orderRpcCli := OrderRpcServCli{
		rpcClient: client,
		timeout:   s.timeout,
	}
	return orderRpcCli
}

func (s *ServClient) GetProductRpcServCli() ProductRpcServCli {
	client := pb.NewProductServClient(s.conn)
	productRpcCli := ProductRpcServCli{
		rpcClient: client,
		timeout:   s.timeout,
	}
	return productRpcCli
}

var _ service.UserServ = (*UserRpcServCli)(nil)

var _ service.OrderServ = (*OrderRpcServCli)(nil)

var _ service.ProductServ = (*ProductRpcServCli)(nil)

func (o *OrderRpcServCli) AddOrder(order *entity.Order) error {
	req := &pb.CreateOrderRequest{
		Order: o.changeFromEntityToRpc(order),
	}
	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	_, err := o.rpcClient.CreateOrder(ctx, req)
	return err
}

func (o *OrderRpcServCli) GetOrderById(id uint) (entity.Order, error) {

	req := &pb.GetOrderByIdRequest{
		Id: uint32(id),
	}

	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	reply, err := o.rpcClient.GetOrderById(ctx, req)

	var order entity.Order

	if err != nil {
		return order, err
	}

	order = o.changeFromRpcToEntity(reply.GetOrder())
	return order, nil

}

func (o *OrderRpcServCli) GetOrdersByUID(uid uint) ([]entity.Order, error) {
	req := &pb.GetOrderByUIdRequest{
		Uid: uint32(uid),
	}

	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	reply, err := o.rpcClient.GetOrderByUId(ctx, req)

	orders := make([]entity.Order, 0)

	if err != nil {
		return orders, err
	}

	orders = o.changeFromRpcToEntitys(reply.GetOrders())

	return orders, nil
}

func (o *OrderRpcServCli) GetOrdersByPID(pid uint) ([]entity.Order, error) {
	req := &pb.GetOrderByPIdRequest{
		Pid: uint32(pid),
	}

	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	reply, err := o.rpcClient.GetOrderByPId(ctx, req)

	orders := make([]entity.Order, 0)

	if err != nil {
		return orders, err
	}

	orders = o.changeFromRpcToEntitys(reply.GetOrders())

	return orders, nil
}

func (o *OrderRpcServCli) DeleteOrder(id uint) error {
	req := &pb.DeleteOrderRequest{
		Id: uint32(id),
	}
	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	_, err := o.rpcClient.DeleteOrder(ctx, req)
	return err
}

func (o *OrderRpcServCli) GetOrders() ([]entity.Order, error) {
	req := &pb.GetOrdersRequest{}

	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	reply, err := o.rpcClient.GetOrders(ctx, req)

	orders := make([]entity.Order, 0)

	if err != nil {
		return orders, err
	}

	orders = o.changeFromRpcToEntitys(reply.GetOrders())

	return orders, nil
}

func (o *OrderRpcServCli) changeFromEntityToRpc(order *entity.Order) *pb.Order {
	return &pb.Order{
		UserId:    uint32(order.UserId),
		ProductId: uint32(order.ProductId),
		Price:     float32(order.Price),
		Num:       int32(order.Num),
		Created:   timestamppb.New(order.Created),
	}
}

func (o *OrderRpcServCli) changeFromRpcToEntity(order *pb.Order) entity.Order {
	return entity.Order{
		ID:        uint(order.GetId()),
		UserId:    uint(order.GetUserId()),
		ProductId: uint(order.GetProductId()),
		Price:     float64(order.GetPrice()),
		Num:       int(order.GetNum()),
		Created:   order.GetCreated().AsTime(),
	}
}

func (o *OrderRpcServCli) changeFromRpcToEntitys(orders []*pb.Order) []entity.Order {
	entityOs := make([]entity.Order, 0)
	for _, order := range orders {
		entityOs = append(entityOs, o.changeFromRpcToEntity(order))
	}
	return entityOs
}

func (u *UserRpcServCli) AddUser(user *entity.User) error {
	req := &pb.CreateUserRequest{
		User: u.changeFromEntityToRpc(user),
	}
	ctx, cancle := context.WithTimeout(context.Background(), u.timeout)
	defer cancle()
	_, err := u.rpcClient.CreateUser(ctx, req)
	return err
}

func (u *UserRpcServCli) GetUser(name string) (entity.User, error) {
	req := &pb.GetUserRequest{
		Name: name,
	}
	var user entity.User
	ctx, cancle := context.WithTimeout(context.Background(), u.timeout)
	defer cancle()

	reply, err := u.rpcClient.GetUser(ctx, req)
	if err != nil {
		return user, err
	}
	user = u.changeFromRpcToEntity(reply.GetUser())
	return user, nil
}

func (u *UserRpcServCli) DeleteUser(name string) error {
	req := &pb.DeleteUserRequest{
		Name: name,
	}
	ctx, cancle := context.WithTimeout(context.Background(), u.timeout)
	defer cancle()

	_, err := u.rpcClient.DeleteUser(ctx, req)
	return err
}

func (u *UserRpcServCli) GetUsers() ([]entity.User, error) {
	req := &pb.GetUsersRequest{}
	users := make([]entity.User, 0)

	ctx, cancle := context.WithTimeout(context.Background(), u.timeout)
	defer cancle()

	reply, err := u.rpcClient.GetUsers(ctx, req)
	if err != nil {
		return users, err
	}
	users = u.changeFromRpcToEntitys(reply.GetUsers())
	return users, nil
}

func (u *UserRpcServCli) changeFromEntityToRpc(user *entity.User) *pb.User {
	return &pb.User{
		Name:    user.Name,
		Sex:     int32(user.Sex),
		Phone:   user.Phone,
		Created: timestamppb.New(user.Created),
	}
}

func (u *UserRpcServCli) changeFromRpcToEntity(user *pb.User) entity.User {
	return entity.User{
		ID:      uint(user.GetId()),
		Name:    user.GetName(),
		Sex:     int(user.GetSex()),
		Phone:   user.GetPhone(),
		Created: user.GetCreated().AsTime(),
	}
}

func (u *UserRpcServCli) changeFromRpcToEntitys(users []*pb.User) []entity.User {
	entityUs := make([]entity.User, 0)
	for _, user := range users {
		entityUs = append(entityUs, u.changeFromRpcToEntity(user))
	}
	return entityUs
}

func (p *ProductRpcServCli) AddProduct(product *entity.Product) error {
	req := &pb.CreateProductRequest{
		Product: p.changeFromEntityToRpc(product),
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	_, err := p.rpcClient.CreateProduct(ctx, req)
	return err
}

func (p *ProductRpcServCli) GetProduct(name string) (entity.Product, error) {
	req := &pb.GetProductRequest{
		Name: name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	reply, err := p.rpcClient.GetProduct(ctx, req)
	var product entity.Product
	if err != nil {
		return product, err
	}
	product = p.changeFromRpcToEntity(reply.GetProduct())
	return product, err

}

func (p *ProductRpcServCli) GetProducts() ([]entity.Product, error) {
	req := &pb.GetProductsRequest{}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	reply, err := p.rpcClient.GetProducts(ctx, req)
	var products []entity.Product
	if err != nil {
		return products, err
	}
	products = p.changeFromRpcToEntitys(reply.GetProducts())
	return products, err
}

func (p *ProductRpcServCli) DeleteProduct(name string) error {
	req := &pb.DeleteProductRequest{
		Name: name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	_, err := p.rpcClient.DeleteProduct(ctx, req)
	return err

}

func (p *ProductRpcServCli) SetStock(id uint, num int) error {
	req := &pb.SetStockRequest{
		Id:  uint32(id),
		Num: int32(num),
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	_, err := p.rpcClient.SetStock(ctx, req)
	return err
}

func (p *ProductRpcServCli) changeFromEntityToRpc(product *entity.Product) *pb.Product {
	return &pb.Product{
		Name:        product.Name,
		Price:       float32(product.Price),
		Stock:       int32(product.Stock),
		Description: product.Description,
		Created:     timestamppb.New(product.Created),
	}
}

func (p *ProductRpcServCli) changeFromRpcToEntity(product *pb.Product) entity.Product {
	return entity.Product{
		ID:          uint(product.GetId()),
		Name:        product.GetName(),
		Price:       float64(product.GetPrice()),
		Stock:       int(product.GetStock()),
		Description: product.GetDescription(),
		Created:     product.GetCreated().AsTime(),
	}
}

func (p *ProductRpcServCli) changeFromRpcToEntitys(products []*pb.Product) []entity.Product {
	entityPs := make([]entity.Product, 0)
	for _, product := range products {
		entityPs = append(entityPs, p.changeFromRpcToEntity(product))
	}
	return entityPs
}