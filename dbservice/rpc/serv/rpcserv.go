package serv

import (
	"context"
	"fmt"
	"net"
	entity2 "seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/dbservice/service"
	registercenter "seckill/registercenter/registerservice"
	pb "seckill/rpc/dbservice"

	"google.golang.org/grpc"
)

type RpcServServer struct {
	pb.UnimplementedUserServServer
	pb.UnimplementedOrderServServer
	pb.UnimplementedProductServServer
	pb.UnimplementedPermServServer
	UserServ       interfaces.UserServ
	ProductServ    interfaces.ProductServ
	OrderServ      interfaces.OrderServ
	PermServ       interfaces.PermServ
	registerCenter *registercenter.RegisterCenter
}

func (s *RpcServServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	user := changeFromURpcToEntity(in.GetUser())
	reply := &pb.CreateUserReply{}
	err := s.UserServ.AddUser(user)
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}

func (s *RpcServServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserReply, error) {
	reply := &pb.GetUserReply{}
	user, err := s.UserServ.GetUser(in.GetName())
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.User = changeFromUEntityToRpc(&user)
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	reply := &pb.DeleteUserReply{}
	err := s.UserServ.DeleteUser(in.GetName())
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) GetUsers(ctx context.Context, in *pb.GetUsersRequest) (*pb.GetUsersReply, error) {
	reply := &pb.GetUsersReply{}
	users, err := s.UserServ.GetUsers()
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Users = changeFromUEntitysToRpc(users)
	reply.Ok = true
	return reply, nil
}

func (s *RpcServServer) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	reply := &pb.CreateOrderReply{}
	order := changeFromORpcToEntity(in.GetOrder())
	err := s.OrderServ.AddOrder(&order, in.GetMethod())
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) GetOrderById(ctx context.Context, in *pb.GetOrderByIdRequest) (*pb.GetOrderByIdReply, error) {
	reply := &pb.GetOrderByIdReply{}
	order, err := s.OrderServ.GetOrderById(uint(in.GetId()))
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Order = changeFromOEntityToRpc(&order)
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) GetOrderByUId(ctx context.Context, in *pb.GetOrderByUIdRequest) (*pb.GetOrderByUIdReply, error) {
	reply := &pb.GetOrderByUIdReply{}
	orders, err := s.OrderServ.GetOrdersByUID(uint(in.GetUid()))
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Orders = changeFromOEntitysToRpc(orders)
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) GetOrderByPId(ctx context.Context, in *pb.GetOrderByPIdRequest) (*pb.GetOrderByPIdReply, error) {
	reply := &pb.GetOrderByPIdReply{}
	orders, err := s.OrderServ.GetOrdersByPID(uint(in.GetPid()))
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Orders = changeFromOEntitysToRpc(orders)
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) DeleteOrder(ctx context.Context, in *pb.DeleteOrderRequest) (*pb.DeleteOrderReply, error) {
	reply := &pb.DeleteOrderReply{}
	err := s.OrderServ.DeleteOrder(uint(in.GetId()))
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) GetOrders(ctx context.Context, in *pb.GetOrdersRequest) (*pb.GetOrdersReply, error) {
	reply := &pb.GetOrdersReply{}
	orders, err := s.OrderServ.GetOrders()
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Orders = changeFromOEntitysToRpc(orders)
	reply.Ok = true
	return reply, nil
}

func (s *RpcServServer) ClearOrders(ctx context.Context, in *pb.ClearOrdersRequest) (*pb.ClearOrdersReply, error) {
	reply := &pb.ClearOrdersReply{}
	err := s.OrderServ.ClearOrders()
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}

func (s *RpcServServer) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.CreateProductReply, error) {
	reply := &pb.CreateProductReply{}
	product := changeFromPRpcToEntity(in.GetProduct())
	err := s.ProductServ.AddProduct(&product)
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.GetProductReply, error) {
	reply := &pb.GetProductReply{}
	product, err := s.ProductServ.GetProduct(in.GetName())
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	p := changeFromPEntityToRpc(&product)
	reply.Product = p
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) DeleteProduct(ctx context.Context, in *pb.DeleteProductRequest) (*pb.DeleteProductReply, error) {
	reply := &pb.DeleteProductReply{}
	err := s.ProductServ.DeleteProduct(in.GetName())
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) GetProducts(ctx context.Context, in *pb.GetProductsRequest) (*pb.GetProductsReply, error) {
	reply := &pb.GetProductsReply{}
	products, err := s.ProductServ.GetProducts()
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Products = changeFromPEntitysToRpc(products)
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) SetStock(ctx context.Context, in *pb.SetStockRequest) (*pb.SetStockReply, error) {
	reply := &pb.SetStockReply{}
	err := s.ProductServ.SetStock(uint(in.GetId()), int(in.GetNum()))
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}

func (s *RpcServServer) GetStock(ctx context.Context, in *pb.GetStockRequest) (*pb.GetStockReply, error) {
	reply := &pb.GetStockReply{}
	stock, err := s.ProductServ.GetStock(uint(in.GetId()), "")
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Stock = int32(stock)
	reply.Ok = true
	return reply, nil
}

func (s *RpcServServer) GetRoles(ctx context.Context, in *pb.GetRolesRequest) (*pb.GetRolesReply, error) {
	reply := &pb.GetRolesReply{}
	roles, err := s.PermServ.RoleList()
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Roles = changeFromRoleEntitysToRpc(roles)
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) GetPerm(ctx context.Context, in *pb.GetPermRequest) (*pb.GetPermReply, error) {
	reply := &pb.GetPermReply{}
	perm, err := s.PermServ.GetPerm(uint(in.GetUid()))
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Perms = perm
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) GetPerms(ctx context.Context, in *pb.GetPermsRequest) (*pb.GetPermsReply, error) {
	reply := &pb.GetPermsReply{}
	perms, err := s.PermServ.PermList()
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Perms = changeFromPermEntitysToRpc(perms)
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) AddRole(ctx context.Context, in *pb.AddRoleRequest) (*pb.AddRoleReply, error) {
	reply := &pb.AddRoleReply{}
	role := changeFromRoleRpcToEntity(in.GetRole())
	err := s.PermServ.AddRole(&role)
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) AddPerm(ctx context.Context, in *pb.AddPermRequest) (*pb.AddPermReply, error) {
	reply := &pb.AddPermReply{}
	perm := changeFromPermRpcToEntity(in.GetPerm())
	err := s.PermServ.AddPerm(&perm)
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}
func (s *RpcServServer) SetRole(ctx context.Context, in *pb.SetRoleRequest) (*pb.SetRoleReply, error) {
	reply := &pb.SetRoleReply{}
	err := s.PermServ.SetRole(int(in.GetUid()), int(in.GetRid()))
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}

func changeFromPEntityToRpc(product *entity2.Product) *pb.Product {
	return &pb.Product{
		Id:          uint32(product.ID),
		Name:        product.Name,
		Price:       float32(product.Price),
		Stock:       int32(product.Stock),
		Description: product.Description,
		Created:     product.Created,
		Version:     int32(product.Version),
	}
}

func changeFromPRpcToEntity(product *pb.Product) entity2.Product {
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

func changeFromPEntitysToRpc(products []entity2.Product) []*pb.Product {
	entityPs := make([]*pb.Product, 0)
	for _, product := range products {
		entityPs = append(entityPs, changeFromPEntityToRpc(&product))
	}
	return entityPs
}

func changeFromUEntityToRpc(user *entity2.User) *pb.User {
	return &pb.User{
		Id:       uint32(user.ID),
		Name:     user.Name,
		Password: user.Password,
		Sex:      int32(user.Sex),
		Phone:    user.Phone,
		Created:  user.Created,
	}
}

func changeFromURpcToEntity(user *pb.User) *entity2.User {
	return &entity2.User{
		ID:       uint(user.GetId()),
		Name:     user.GetName(),
		Password: user.GetPassword(),
		Sex:      int(user.GetSex()),
		Phone:    user.GetPhone(),
		Created:  user.GetCreated(),
	}
}

func changeFromUEntitysToRpc(users []entity2.User) []*pb.User {
	entityUs := make([]*pb.User, 0)
	for _, user := range users {
		entityUs = append(entityUs, changeFromUEntityToRpc(&user))
	}
	return entityUs
}

func changeFromOEntityToRpc(order *entity2.Order) *pb.Order {
	return &pb.Order{
		Id:        uint32(order.ID),
		UserId:    uint32(order.UserId),
		ProductId: uint32(order.ProductId),
		Price:     float32(order.Price),
		Num:       int32(order.Num),
		Created:   order.Created,
	}
}

func changeFromORpcToEntity(order *pb.Order) entity2.Order {
	return entity2.Order{
		ID:        uint(order.GetId()),
		UserId:    uint(order.GetUserId()),
		ProductId: uint(order.GetProductId()),
		Price:     float64(order.GetPrice()),
		Num:       int(order.GetNum()),
		Created:   order.GetCreated(),
	}
}

func changeFromOEntitysToRpc(orders []entity2.Order) []*pb.Order {
	entityOs := make([]*pb.Order, 0)
	for _, order := range orders {
		entityOs = append(entityOs, changeFromOEntityToRpc(&order))
	}
	return entityOs
}

func changeFromRoleEntitysToRpc(roles []entity2.Role) []*pb.Role {
	entityRs := make([]*pb.Role, 0)
	for _, role := range roles {
		entityRs = append(entityRs, changeFromRoleEntityToRpc(&role))
	}
	return entityRs
}

func changeFromRoleEntityToRpc(role *entity2.Role) *pb.Role {
	return &pb.Role{
		Id:   uint32(role.ID),
		Name: role.Name,
	}
}

func changeFromRoleRpcToEntity(role *pb.Role) entity2.Role {
	return entity2.Role{
		ID:   uint(role.GetId()),
		Name: role.GetName(),
	}
}

func changeFromPermEntitysToRpc(perms []entity2.Perm) []*pb.Perm {
	entityPs := make([]*pb.Perm, 0)
	for _, perm := range perms {
		entityPs = append(entityPs, changeFromPermEntityToRpc(&perm))
	}
	return entityPs
}

func changeFromPermEntityToRpc(perm *entity2.Perm) *pb.Perm {
	return &pb.Perm{
		Id:   uint32(perm.ID),
		Path: perm.Path,
	}
}

func changeFromPermRpcToEntity(perm *pb.Perm) entity2.Perm {
	return entity2.Perm{
		ID:   uint(perm.GetId()),
		Path: perm.GetPath(),
	}
}

func (s *RpcServServer) StartRpcServServer(port string) error {
	// setting, err := seetings.GetSetting()
	// if err != nil {
	// 	fmt.Println("get settings error in start rpc serv!")
	// 	return err
	// }
	//port := setting.RPC.DbServPort
	lis, err := net.Listen("tcp", port)
	s.registerCenter.Register("/bwsk/dbservice/"+port, port)
	if err != nil {
		fmt.Printf("rpc service listen port %s error", port)
		return err
	}
	s.initServ()
	serv := grpc.NewServer()
	pb.RegisterUserServServer(serv, s)
	pb.RegisterOrderServServer(serv, s)
	pb.RegisterProductServServer(serv, s)
	pb.RegisterPermServServer(serv, s)
	if err = serv.Serve(lis); err != nil {
		fmt.Println("serv rpc serve error")
		return err
	}
	return nil
}

func (s *RpcServServer) initServ() {
	if s.UserServ == nil {
		s.UserServ = service.GetUserServ()
	}
	if s.ProductServ == nil {
		s.ProductServ = service.GetProductServ()
	}
	if s.OrderServ == nil {
		s.OrderServ = service.GetOrderServ()
	}
	if s.PermServ == nil {
		s.PermServ = service.GetPermServ()
	}
}

func NewRpcServServer() *RpcServServer {

	rpcServ := &RpcServServer{
		UserServ:       service.GetUserServ(),
		ProductServ:    service.GetProductServ(),
		OrderServ:      service.GetOrderServ(),
		PermServ:       service.GetPermServ(),
		registerCenter: registercenter.GetRegisterCenter(),
	}

	return rpcServ
}
