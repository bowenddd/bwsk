package service

import (
	"errors"
	"fmt"
	"seckill/cacheservice/store"
	"seckill/common/consts"
	"seckill/common/entity"
	"seckill/common/interfaces"
	dbrpc "seckill/dbservice/rpc"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheServImpl struct {
	store *store.DataStore
}

func (c *CacheServImpl) GetStock(id uint) (int, error) {
	id_str := strconv.Itoa(int(id))
	stock_str, err := c.store.Get(id_str)
	if err == redis.Nil {
		// 如果在redis中找不到的话，需要从数据库中查找，然后将数据写入redis
		serv, err := dbrpc.GetDbServRpcCli()
		if err != nil {
			return 0, err
		}
		prsc := serv.GetProductRpcServCli()
		stock, err := prsc.GetStock(id, "")
		if err != nil {
			return 0, nil
		}
		err = c.SetStock(id, stock, 0)
		if err != nil {
			return 0, err
		}
		return stock, nil
	}
	if err != nil {
		return 0, err
	}
	stock, err := strconv.Atoi(stock_str)
	if err != nil {
		return 0, err
	}
	return stock, nil
}

func (c *CacheServImpl) SetStock(id uint, num int, exp time.Duration) error {
	id_str := strconv.Itoa(int(id))
	num_str := strconv.Itoa(num)
	err := c.store.Set(id_str, num_str, exp)
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheServImpl) CreateOrder(order *entity.Order, method string) error {
	switch method {
	case consts.CACHEOPTIMISTICLOCK:
		return c.store.CreateOrderByOLOCK(order)
	case consts.CACHEPESSIMISTICLOCK:
		return c.store.CreateOrderByPLOCK(order)
	default:
		return errors.New("method not support")
	}
}

func (c *CacheServImpl) Lock(key string, ex time.Duration) (bool, error) {
	return c.store.Lock(key, ex)
}

func (c *CacheServImpl) UnLock(key string) int64 {
	return c.store.UnLock(key)
}

func (c *CacheServImpl) GetUserPerms(id uint) (string, error) {
	key := fmt.Sprintf("user_id_%d_perms", id)
	perms, err := c.store.Get(key)
	if errors.Is(err, redis.Nil) {
		rpcCli, rpcErr := dbrpc.GetDbServRpcCli()
		if rpcErr != nil {
			return "", rpcErr
		}
		prsc := rpcCli.GetPermRpcServCli()
		perms, rpcErr = prsc.GetPerm(id)
		if rpcErr != nil {
			return "", rpcErr
		}
		storeErr := c.store.Set(key, perms, 0)
		if storeErr != nil {
			return "", storeErr
		}
		return perms, nil
	} else if err != nil {
		return "", err
	}
	return perms, nil
}

var cacheServImpl *CacheServImpl

var once = new(sync.Once)

func GetCacheServ() interfaces.CacheServ {
	once.Do(func() {
		cacheServImpl = &CacheServImpl{store: store.GetDataStore()}
	})
	return cacheServImpl
}
