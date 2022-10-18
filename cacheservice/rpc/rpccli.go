package rpc

import (
	"fmt"
	"math"
	"seckill/common/interfaces"
	pb "seckill/rpc/cacheservice"
	"time"

	"seckill/seetings"

	"seckill/common/entity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"context"
)

type CacheServCli struct {
	client  pb.CacheServiceClient
	timeout time.Duration
}

var _ interfaces.CacheServ = (*CacheServCli)(nil)

func (c *CacheServCli) SetStock(id uint, num int, exp time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	_, err := c.client.SetStock(ctx, &pb.SetStockRequest{
		ProductId: uint32(id),
		Num:       int32(num),
		Expire:    int32(exp),
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *CacheServCli) GetStock(id uint) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.client.GetStock(ctx, &pb.GetStockRequest{
		ProductId: uint32(id),
	})
	if err != nil {
		return 0, err
	}
	return int(resp.Num), nil
}
func (c *CacheServCli) CreateOrder(order *entity.Order, method string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	request := &pb.CreateOrderRequest{
		Order: &pb.Order{
			Id:        uint32(order.ID),
			ProductId: uint32(order.ProductId),
			UserId:    uint32(order.UserId),
			Created:   order.Created,
			Num:       int32(order.Num),
			Price:     float32(order.Price),
		},
		Method: method,
	}
	_, err := c.client.CreateOrder(ctx, request)
	if err != nil {
		return err
	}
	return nil
}
func (c *CacheServCli) Lock(key string, ex time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.client.Lock(ctx, &pb.LockRequest{
		Key:    key,
		Expire: int32(ex),
	})
	if err != nil {
		return false, err
	}
	return resp.Ok, nil
}
func (c *CacheServCli) UnLock(key string) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.client.Unlock(ctx, &pb.UnlockRequest{Key: key})
	if err != nil {
		return -1
	}
	return int64(resp.N)
}

func (c *CacheServCli) GetUserPerms(userId uint) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.client.GetUserPerms(ctx, &pb.GetUserPermsRequest{
		UserId: uint32(userId),
	})
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return resp.Perms, nil
}

func CreateCacheServClient(addr string, timeout time.Duration) (*CacheServCli, error) {
	var conn *grpc.ClientConn
	var err error
	if timeout == -1 {
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		timeout = math.MaxInt64
	} else {
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(timeout))
	}
	if err != nil {
		fmt.Println("cacheservice:grpc dial error! in new serv client")
		return &CacheServCli{}, err
	}
	client := pb.NewCacheServiceClient(conn)
	return &CacheServCli{
		client:  client,
		timeout: timeout,
	}, nil
}

func getRpcSeetings() (addr string, timeout int, err error) {
	seeting, err := seetings.GetSetting()
	if err != nil {
		fmt.Println("cacheservice: get seeting error! in get rpc seeting")
		return
	}
	port := seeting.RPC.CacheServPort
	addr = fmt.Sprintf("localhost%s", port)
	timeout = seeting.RPC.Timeout
	return
}

func NewCacheServClient() (*CacheServCli, error) {
	addr, timeout, err := getRpcSeetings()
	var cli *CacheServCli
	if err != nil {
		fmt.Println("cacheservice: get rpc seeting error! in get cache serv client")
		return cli, err
	}
	cli, err = CreateCacheServClient(addr, time.Duration(timeout)*time.Second)
	if err != nil {
		fmt.Println("cacheservice: new cache serv client error! in get cache serv client")
		return cli, err
	}
	return cli, nil
}
