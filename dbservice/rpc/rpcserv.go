package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"seckill/dbservice/service"
	"seckill/entity"
	pb "seckill/rpc/dbservice"
	"seckill/seetings"
	"sync"
)

type RpcServServer struct {
	pb.UnimplementedUserServServer
	pb.UnimplementedOrderServServer
	pb.UnimplementedProductServServer
	UserServ    service.UserServ
	ProductServ service.ProductServ
	OrderServ   service.OrderServ
}

func (s *RpcServServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	user := &entity.User{
		Name:    in.GetUser().GetName(),
		Sex:     int(in.GetUser().GetSex()),
		Phone:   in.GetUser().GetPhone(),
		Created: in.GetUser().GetCreated().AsTime(),
	}
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
	reply.User = &pb.User{
		Id:      uint32(user.ID),
		Name:    user.Name,
		Sex:     int32(user.Sex),
		Phone:   user.Phone,
		Created: timestamppb.New(user.Created),
	}
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
	for _, user := range users {
		u := &pb.User{
			Id:      uint32(user.ID),
			Name:    user.Name,
			Phone:   user.Phone,
			Sex:     int32(user.Sex),
			Created: timestamppb.New(user.Created),
		}
		reply.Users = append(reply.Users, u)
	}
	reply.Ok = true
	return reply, nil
}

func (s *RpcServServer) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	reply := &pb.CreateOrderReply{}
	order := &entity.Order{
		UserId:    uint(in.GetOrder().GetUserId()),
		ProductId: uint(in.GetOrder().GetProductId()),
		Num:       int(in.GetOrder().GetNum()),
		Price:     float64(in.GetOrder().GetPrice()),
		Created:   in.GetOrder().GetCreated().AsTime(),
	}
	err := s.OrderServ.AddOrder(order)
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
	o := &pb.Order{
		Id:        uint32(order.ID),
		UserId:    uint32(order.UserId),
		ProductId: uint32(order.ProductId),
		Num:       int32(order.Num),
		Price:     float32(order.Price),
		Created:   timestamppb.New(order.Created),
	}
	reply.Order = o
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
	for _, order := range orders {
		o := &pb.Order{
			Id:        uint32(order.ID),
			UserId:    uint32(order.UserId),
			ProductId: uint32(order.ProductId),
			Num:       int32(order.Num),
			Price:     float32(order.Price),
			Created:   timestamppb.New(order.Created),
		}
		reply.Orders = append(reply.Orders, o)
	}
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
	for _, order := range orders {
		o := &pb.Order{
			Id:        uint32(order.ID),
			UserId:    uint32(order.UserId),
			ProductId: uint32(order.ProductId),
			Num:       int32(order.Num),
			Price:     float32(order.Price),
			Created:   timestamppb.New(order.Created),
		}
		reply.Orders = append(reply.Orders, o)
	}
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
	for _, order := range orders {
		o := &pb.Order{
			Id:        uint32(order.ID),
			UserId:    uint32(order.UserId),
			ProductId: uint32(order.ProductId),
			Num:       int32(order.Num),
			Price:     float32(order.Price),
			Created:   timestamppb.New(order.Created),
		}
		reply.Orders = append(reply.Orders, o)
	}
	reply.Ok = true
	return reply, nil
}

func (s *RpcServServer) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.CreateProductReply, error) {
	reply := &pb.CreateProductReply{}
	product := &entity.Product{
		Name:        in.GetProduct().GetName(),
		Price:       float64(in.GetProduct().GetPrice()),
		Stock:       int(in.GetProduct().GetStock()),
		Description: in.GetProduct().GetDescription(),
		Created:     in.GetProduct().GetCreated().AsTime(),
	}
	err := s.ProductServ.AddProduct(product)
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
	p := &pb.Product{
		Id:          uint32(product.ID),
		Name:        product.Name,
		Price:       float32(product.Price),
		Stock:       int32(product.Stock),
		Description: product.Description,
		Created:     timestamppb.New(product.Created),
	}
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
	for _, product := range products {
		p := &pb.Product{
			Id:          uint32(product.ID),
			Name:        product.Name,
			Price:       float32(product.Price),
			Stock:       int32(product.Stock),
			Description: product.Description,
			Created:     timestamppb.New(product.Created),
		}
		reply.Products = append(reply.Products, p)
	}
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

func (s *RpcServServer) StartRpcServServer() error {
	setting, err := seetings.GetSetting()
	if err != nil {
		fmt.Println("get settings error in start rpc serv!")
		return err
	}
	port := setting.RPC.DbServPort
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("rpc service listen port %s error", port)
		return err
	}
	s.initServ()
	serv := grpc.NewServer()
	pb.RegisterUserServServer(serv, s)
	pb.RegisterOrderServServer(serv, s)
	pb.RegisterProductServServer(serv, s)
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
}

var rpcServ *RpcServServer

var mu sync.Mutex

func init() {
	mu = sync.Mutex{}
}

func GetRpcServServer() *RpcServServer {
	if rpcServ == nil {
		mu.Lock()
		if rpcServ == nil {
			rpcServ = &RpcServServer{
				UserServ:    service.GetUserServ(),
				ProductServ: service.GetProductServ(),
				OrderServ:   service.GetOrderServ(),
			}
		}
		mu.Unlock()
	}
	return rpcServ
}
