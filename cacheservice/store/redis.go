package store

import (
	"context"
	"fmt"
	"seckill/common/consts"
	"seckill/common/entity"
	"seckill/dbservice/rpc"
	"seckill/seetings"
	"strconv"
	"time"
	"errors"
	"github.com/go-redis/redis/v8"
	"sync"
)

type DataStore struct {
	client *redis.Client
}

var redisStore *DataStore

var once = new(sync.Once)

func init(){

}

func GetDataStore() *DataStore{
	once.Do(func() {
		setting, err := seetings.GetSetting()
		if err != nil {
			panic(err)
		}
		addr := fmt.Sprintf("%s:%d", setting.Redis.Host, setting.Redis.Port)
		redisStore = &DataStore{
			client: redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: setting.Redis.Password,
				DB:       setting.Redis.Db,
			}),
		}
	})
	return redisStore
}

func (ds *DataStore) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return ds.client.Set(ctx, key, value, expiration).Err()
}

func (ds *DataStore) Get(key string) (string, error) {
	ctx := context.Background()
	return ds.client.Get(ctx, key).Result()
}

// 加锁
func (ds *DataStore)Lock(key string,ex time.Duration)(bool,error) {
	// ex:设置默认过期时间10秒，防止死锁
	ctx := context.Background()
	value := "lock_for_key_" + key
	succ := false
	var err error
	timeout := false
	st := time.Now()
	for !succ && !timeout{
		succ ,err = ds.client.SetNX(ctx,key, value, ex).Result()
		if err != nil{
			return false, err
		}
		if time.Since(st) > ex{
			timeout = true
		}
	}
	if timeout{
		return false, errors.New("获取锁超时")
	}
	return succ, nil
}

// 解锁
func (ds *DataStore)UnLock(key string) int64 {
	ctx := context.Background()
	nums, err := ds.client.Del(ctx,key).Result()
	for err != nil {
		nums, err = ds.client.Del(ctx,key).Result()
	}
	return nums
}


// redis悲观锁解决缓存层超卖问题
func (ds *DataStore) CreateOrderByPLOCK(order *entity.Order) error{
	ctx := context.Background()
	pid := strconv.Itoa(int(order.ProductId))
	rpcCli, err := rpc.GetDbServRpcCli()
	if err != nil{
		return err
	}
	// 加锁
	_, err = ds.Lock(pid,10*time.Second)
	defer ds.UnLock(pid)
	if err != nil{
		return err
	}
	// 获取库存
	stock, err := ds.client.Get(ctx, pid).Int()
	if err != nil{
		return err
	}
	if stock < order.Num{
		return errors.New("库存不足")
	}
	sc := ds.client.Set(ctx, pid, stock-order.Num, 0)
	if sc.Err() != nil{
		return sc.Err()
	}
	orderServCli := rpcCli.GetOrderRpcServCli()
	err = orderServCli.AddOrder(order, consts.NOMEASURE)
	if err != nil{
		sc := ds.client.Set(ctx, pid, stock, 0)
		for sc.Err()!= nil{
			sc = ds.client.Set(ctx, pid, stock, 0)
		}
		return err
	}
	return nil
}

// redis乐观锁--watch+transaction解决缓存层超卖问题
func (ds *DataStore) CreateOrderByOLOCK(order *entity.Order) error{
	ctx := context.Background()

	txf := func (tx *redis.Tx) error{
		return ds.CreateOrderTransaction(tx, order)
	}
	pid := strconv.Itoa(int(order.ProductId))
	for i := 0; i < consts.REDIS_CREATE_ORDER_MAX_RETRY; i++ {
		err := ds.client.Watch(ctx, txf, pid)
		if err == nil {
			return nil
		}
		if err == redis.TxFailedErr {
			continue
		}
		return err
	}
	return errors.New("increment reached maximum number of retries")
}

func (ds *DataStore)CreateOrderTransaction(tx *redis.Tx,order *entity.Order) error{
	ctx := context.Background()
	pid := strconv.Itoa(int(order.ProductId))
	stock, err := tx.Get(ctx, pid).Int()
	if stock < order.Num{
		return errors.New("库存不足")
	}
	if err != nil {
		return err
	}

	rpcCli, err := rpc.GetDbServRpcCli()
	if err != nil{
		return err
	}

	// 缓存减少库存
	tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		sc := pipe.Set(ctx, pid, stock-order.Num, 0)
		return sc.Err()
	})

	// 通过RPC调用数据库服务，创建订单
	orderServCli := rpcCli.GetOrderRpcServCli()
	err = orderServCli.AddOrder(order, consts.NOMEASURE)
	if err != nil{
		// 回滚
		incrErr := ds.increment(pid)
		if incrErr != nil{
			return incrErr
		}
		return err
	}
	return nil
}

func (ds *DataStore)increment(key string) error {
	ctx := context.Background()

	txf := func(tx *redis.Tx) error {

		n, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}
	
		n++

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, n, 0)
			return nil
		})

		return err
	}

	for i := 0; i < 1000; i++ {
		err := ds.client.Watch(ctx, txf, key)
		if err == nil {
			return nil
		}
		if err == redis.TxFailedErr {
			continue
		}
		return err
	}

	return errors.New("increment reached maximum number of retries")
}
