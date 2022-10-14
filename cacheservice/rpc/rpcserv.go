package rpc

import (
	"context"
	"net"
	"seckill/cacheservice/service"
	"seckill/common/entity"
	"seckill/common/interfaces"
	pb "seckill/rpc/cacheservice"
	"seckill/seetings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type CacheRpcServ struct {
	pb.UnimplementedCacheServiceServer
	serv interfaces.CacheServ
}

func (s *CacheRpcServ) initService() {
	if s.serv == nil {
		s.serv = service.GetCacheServ()
	}
}

func (s *CacheRpcServ) StartCacheRpcServService() {
	seeting, err := seetings.GetSetting()
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", seeting.RPC.CacheServPort)
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

func ChangeFromRpcToEntity(order *pb.Order) *entity.Order {
	return &entity.Order{
		ProductId: uint(order.GetProductId()),
		UserId:    uint(order.GetUserId()),
		Num:       int(order.GetNum()),
		Price:     float64(order.GetPrice()),
		Created:   order.GetCreated(),
	}
}

var cacheRpcServService *CacheRpcServ

var cacheRpcServOnce = new(sync.Once)


func GetCacheRpcService() *CacheRpcServ {
	cacheRpcServOnce.Do(func() {
		cacheRpcServService = &CacheRpcServ{
			serv: service.GetCacheServ(),
		}
	})
	return cacheRpcServService
}