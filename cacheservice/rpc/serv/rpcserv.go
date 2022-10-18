package serv

import (
	"context"
	"net"
	"seckill/cacheservice/service"
	"seckill/common/entity"
	"seckill/common/interfaces"
	pb "seckill/rpc/cacheservice"
	"time"
	registercenter "seckill/registercenter/registerservice"
	"google.golang.org/grpc"
)

type CacheRpcServ struct {
	pb.UnimplementedCacheServiceServer
	serv interfaces.CacheServ
	registerCenter *registercenter.RegisterCenter
}

func (s *CacheRpcServ) initService() {
	if s.serv == nil {
		s.serv = service.GetCacheServ()
	}
}

func (s *CacheRpcServ) StartCacheRpcServService(port string) {
	lis, err := net.Listen("tcp", port)
	s.registerCenter.Register("/bwsk/cacheservice/"+port, port)
	if err != nil {
		panic(err)
	}
	s.initService()
	serv := grpc.NewServer()
	pb.RegisterCacheServiceServer(serv, s)
	if err = serv.Serve(lis); err != nil {
		panic(err)
	}
}

func (s *CacheRpcServ) SetStock(ctx context.Context, in *pb.SetStockRequest) (*pb.SetStockReply, error) {
	reply := &pb.SetStockReply{}
	err := s.serv.SetStock(uint(in.GetProductId()), int(in.GetNum()), time.Duration(in.GetExpire())*time.Millisecond)
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil

}
func (s *CacheRpcServ) GetStock(ctx context.Context, in *pb.GetStockRequest) (*pb.GetStockReply, error) {
	reply := &pb.GetStockReply{}
	stock, err := s.serv.GetStock(uint(in.GetProductId()))
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	reply.Num = int32(stock)
	return reply, nil
}
func (s *CacheRpcServ) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	reply := &pb.CreateOrderReply{}
	order := ChangeFromRpcToEntity(in.GetOrder())
	err := s.serv.CreateOrder(order, in.GetMethod())
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	return reply, nil
}
func (s *CacheRpcServ) Lock(ctx context.Context, in *pb.LockRequest) (*pb.LockReply, error) {
	reply := &pb.LockReply{}
	succ, err := s.serv.Lock(in.GetKey(), time.Duration(in.GetExpire())*time.Millisecond)
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = succ
	return reply, nil
}
func (s *CacheRpcServ) Unlock(ctx context.Context, in *pb.UnlockRequest) (*pb.UnlockReply, error) {
	reply := &pb.UnlockReply{}
	num := s.serv.UnLock(in.GetKey())
	reply.N = int32(num)
	return reply, nil
}

func (s *CacheRpcServ) GetUserPerms(ctx context.Context, in *pb.GetUserPermsRequest) (*pb.GetUserPermsReply, error) {
	reply := &pb.GetUserPermsReply{}
	perms, err := s.serv.GetUserPerms(uint(in.GetUserId()))
	if err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return reply, err
	}
	reply.Ok = true
	reply.Perms = perms
	return reply, nil
}

func ChangeFromRpcToEntity(order *pb.Order) *entity.Order {
	return &entity.Order{
		ProductId: uint(order.GetProductId()),
		UserId:    uint(order.GetUserId()),
		Num:       int(order.GetNum()),
		Price:     float64(order.GetPrice()),
		Created:   order.GetCreated(),
	}
}


func NewCacheRpcService() *CacheRpcServ {
	cacheRpcServService := &CacheRpcServ{
		serv: service.GetCacheServ(),
		registerCenter: registercenter.GetRegisterCenter(),
	}
	return cacheRpcServService
}
