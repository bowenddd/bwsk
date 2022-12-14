package rpc

import (
	"context"
	"fmt"
	"math"
	entity2 "seckill/common/entity"
	"seckill/common/interfaces"
	pb "seckill/rpc/dbservice"
	"seckill/seetings"
	"sync"
	"time"

	"google.golang.org/grpc"
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

type PermRpcServCli struct {
	rpcClient pb.PermServClient
	timeout   time.Duration
}

func NewServClient(addr string, timeout time.Duration) (ServClient, error) {
	var conn *grpc.ClientConn
	var err error
	if timeout == -1 {
		conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
		timeout = math.MaxInt64
	} else {
		conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(timeout))
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
	fmt.Printf("client is %v\n", client)
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

func (s *ServClient) GetPermRpcServCli() PermRpcServCli {
	client := pb.NewPermServClient(s.conn)
	permRpcCli := PermRpcServCli{
		rpcClient: client,
		timeout:   s.timeout,
	}
	return permRpcCli
}

var _ interfaces.UserServ = (*UserRpcServCli)(nil)

var _ interfaces.OrderServ = (*OrderRpcServCli)(nil)

var _ interfaces.ProductServ = (*ProductRpcServCli)(nil)

var _ interfaces.PermServ = (*PermRpcServCli)(nil)

func (o *OrderRpcServCli) AddOrder(order *entity2.Order, method string) error {
	req := &pb.CreateOrderRequest{
		Method: method,
		Order:  o.changeFromEntityToRpc(order),
	}
	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	_, err := o.rpcClient.CreateOrder(ctx, req)
	return err
}

func (o *OrderRpcServCli) GetOrderById(id uint) (entity2.Order, error) {

	req := &pb.GetOrderByIdRequest{
		Id: uint32(id),
	}

	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	reply, err := o.rpcClient.GetOrderById(ctx, req)

	var order entity2.Order

	if err != nil {
		return order, err
	}

	order = o.changeFromRpcToEntity(reply.GetOrder())
	return order, nil

}

func (o *OrderRpcServCli) GetOrdersByUID(uid uint) ([]entity2.Order, error) {
	req := &pb.GetOrderByUIdRequest{
		Uid: uint32(uid),
	}

	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	reply, err := o.rpcClient.GetOrderByUId(ctx, req)

	orders := make([]entity2.Order, 0)

	if err != nil {
		return orders, err
	}

	orders = o.changeFromRpcToEntitys(reply.GetOrders())

	return orders, nil
}

func (o *OrderRpcServCli) GetOrdersByPID(pid uint) ([]entity2.Order, error) {
	req := &pb.GetOrderByPIdRequest{
		Pid: uint32(pid),
	}

	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	reply, err := o.rpcClient.GetOrderByPId(ctx, req)

	orders := make([]entity2.Order, 0)

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

func (o *OrderRpcServCli) GetOrders() ([]entity2.Order, error) {
	req := &pb.GetOrdersRequest{}

	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	reply, err := o.rpcClient.GetOrders(ctx, req)

	orders := make([]entity2.Order, 0)

	if err != nil {
		return orders, err
	}

	orders = o.changeFromRpcToEntitys(reply.GetOrders())

	return orders, nil
}

func (o *OrderRpcServCli) ClearOrders() error {
	req := &pb.ClearOrdersRequest{}
	ctx, cancle := context.WithTimeout(context.Background(), o.timeout)
	defer cancle()

	_, err := o.rpcClient.ClearOrders(ctx, req)
	return err
}

func (o *OrderRpcServCli) changeFromEntityToRpc(order *entity2.Order) *pb.Order {
	return &pb.Order{
		UserId:    uint32(order.UserId),
		ProductId: uint32(order.ProductId),
		Price:     float32(order.Price),
		Num:       int32(order.Num),
		Created:   order.Created,
	}
}

func (o *OrderRpcServCli) changeFromRpcToEntity(order *pb.Order) entity2.Order {
	return entity2.Order{
		ID:        uint(order.GetId()),
		UserId:    uint(order.GetUserId()),
		ProductId: uint(order.GetProductId()),
		Price:     float64(order.GetPrice()),
		Num:       int(order.GetNum()),
		Created:   order.GetCreated(),
	}
}

func (o *OrderRpcServCli) changeFromRpcToEntitys(orders []*pb.Order) []entity2.Order {
	entityOs := make([]entity2.Order, 0)
	for _, order := range orders {
		entityOs = append(entityOs, o.changeFromRpcToEntity(order))
	}
	return entityOs
}

func (u *UserRpcServCli) AddUser(user *entity2.User) error {
	req := &pb.CreateUserRequest{
		User: u.changeFromEntityToRpc(user),
	}
	ctx, cancle := context.WithTimeout(context.Background(), u.timeout)
	defer cancle()
	_, err := u.rpcClient.CreateUser(ctx, req)
	return err
}

func (u *UserRpcServCli) GetUser(name string) (entity2.User, error) {
	req := &pb.GetUserRequest{
		Name: name,
	}
	var user entity2.User
	fmt.Println("here")
	fmt.Println(u)
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

func (u *UserRpcServCli) GetUsers() ([]entity2.User, error) {
	req := &pb.GetUsersRequest{}
	users := make([]entity2.User, 0)

	ctx, cancle := context.WithTimeout(context.Background(), u.timeout)
	defer cancle()

	reply, err := u.rpcClient.GetUsers(ctx, req)
	if err != nil {
		return users, err
	}
	users = u.changeFromRpcToEntitys(reply.GetUsers())
	return users, nil
}

func (u *UserRpcServCli) changeFromEntityToRpc(user *entity2.User) *pb.User {
	return &pb.User{
		Name:     user.Name,
		Password: user.Password,
		Sex:      int32(user.Sex),
		Phone:    user.Phone,
		Created:  user.Created,
	}
}

func (u *UserRpcServCli) changeFromRpcToEntity(user *pb.User) entity2.User {
	return entity2.User{
		ID:       uint(user.GetId()),
		Name:     user.GetName(),
		Password: user.GetPassword(),
		Sex:      int(user.GetSex()),
		Phone:    user.GetPhone(),
		Created:  user.GetCreated(),
	}
}

func (u *UserRpcServCli) changeFromRpcToEntitys(users []*pb.User) []entity2.User {
	entityUs := make([]entity2.User, 0)
	for _, user := range users {
		entityUs = append(entityUs, u.changeFromRpcToEntity(user))
	}
	return entityUs
}

func (p *ProductRpcServCli) AddProduct(product *entity2.Product) error {
	req := &pb.CreateProductRequest{
		Product: p.changeFromEntityToRpc(product),
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	_, err := p.rpcClient.CreateProduct(ctx, req)
	return err
}

func (p *ProductRpcServCli) GetProduct(name string) (entity2.Product, error) {
	req := &pb.GetProductRequest{
		Name: name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	reply, err := p.rpcClient.GetProduct(ctx, req)
	var product entity2.Product
	if err != nil {
		return product, err
	}
	product = p.changeFromRpcToEntity(reply.GetProduct())
	return product, err

}

func (p *ProductRpcServCli) GetProducts() ([]entity2.Product, error) {
	req := &pb.GetProductsRequest{}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	reply, err := p.rpcClient.GetProducts(ctx, req)
	var products []entity2.Product
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

func (p *ProductRpcServCli) GetStock(id uint, method string) (int, error) {
	req := &pb.GetStockRequest{
		Id: uint32(id),
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	reply, err := p.rpcClient.GetStock(ctx, req)
	if err != nil {
		return 0, err
	}
	return int(reply.GetStock()), nil
}

func (p *PermRpcServCli) RoleList() ([]entity2.Role, error) {
	req := &pb.GetRolesRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	reply, err := p.rpcClient.GetRoles(ctx, req)
	var roles []entity2.Role
	if err != nil {
		return roles, err
	}
	roles = p.changeFromRoleRpcToEntitys(reply.GetRoles())
	return roles, err
}
func (p *PermRpcServCli) PermList() ([]entity2.Perm, error) {
	req := &pb.GetPermsRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	reply, err := p.rpcClient.GetPerms(ctx, req)
	var perms []entity2.Perm
	if err != nil {
		return perms, err
	}
	perms = p.changeFromPermRpcToEntitys(reply.GetPerms())
	return perms, err
}
func (p *PermRpcServCli) AddRole(*entity2.Role) error {
	req := &pb.AddRoleRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	_, err := p.rpcClient.AddRole(ctx, req)
	return err
}
func (p *PermRpcServCli) AddPerm(*entity2.Perm) error {
	req := &pb.AddPermRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	_, err := p.rpcClient.AddPerm(ctx, req)
	return err
}
func (p *PermRpcServCli) GetPerm(uid uint) (string, error) {
	req := &pb.GetPermRequest{
		Uid: uint32(uid),
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	reply, err := p.rpcClient.GetPerm(ctx, req)
	if err != nil {
		return "", err
	}
	return reply.GetPerms(), nil
}
func (p *PermRpcServCli) SetRole(userId, roleId int) error {
	req := &pb.SetRoleRequest{
		Uid: uint32(userId),
		Rid: uint32(roleId),
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	_, err := p.rpcClient.SetRole(ctx, req)
	return err
}

func (p *ProductRpcServCli) changeFromEntityToRpc(product *entity2.Product) *pb.Product {
	return &pb.Product{
		Name:        product.Name,
		Price:       float32(product.Price),
		Stock:       int32(product.Stock),
		Description: product.Description,
		Created:     product.Created,
		Version:     int32(product.Version),
	}
}

func (p *ProductRpcServCli) changeFromRpcToEntity(product *pb.Product) entity2.Product {
	return entity2.Product{
		ID:          uint(product.GetId()),
		Name:        product.GetName(),
		Price:       float64(product.GetPrice()),
		Stock:       int(product.GetStock()),
		Description: product.GetDescription(),
		Created:     product.GetCreated(),
		Version:     int(product.GetVersion()),
	}
}

func (p *ProductRpcServCli) changeFromRpcToEntitys(products []*pb.Product) []entity2.Product {
	entityPs := make([]entity2.Product, 0)
	for _, product := range products {
		entityPs = append(entityPs, p.changeFromRpcToEntity(product))
	}
	return entityPs
}

func (p *PermRpcServCli) changeFromRoleRpcToEntity(role *pb.Role) entity2.Role {
	return entity2.Role{
		ID:   uint(role.GetId()),
		Name: role.GetName(),
	}
}

func (p *PermRpcServCli) changeFromPermRpcToEntity(perm *pb.Perm) entity2.Perm {
	return entity2.Perm{
		ID:   uint(perm.GetId()),
		Path: perm.GetPath(),
	}
}

func (p *PermRpcServCli) changeFromRoleRpcToEntitys(roles []*pb.Role) []entity2.Role {
	entityPs := make([]entity2.Role, 0)
	for _, role := range roles {
		entityPs = append(entityPs, p.changeFromRoleRpcToEntity(role))
	}
	return entityPs
}

func (p *PermRpcServCli) changeFromPermRpcToEntitys(perms []*pb.Perm) []entity2.Perm {
	entityPs := make([]entity2.Perm, 0)
	for _, perm := range perms {
		entityPs = append(entityPs, p.changeFromPermRpcToEntity(perm))
	}
	return entityPs
}

var dbServRpcCliOnce = new(sync.Once)

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

func NewDbServRpcCli(port string) (*ServClient, error) {
	_, timeout, err := getRpcSettings()
	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf("localhost%s", port)
	cli, err := NewServClient(addr, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	return &cli, nil
}
