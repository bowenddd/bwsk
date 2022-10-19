package registerservice

import (
	"context"
	"errors"
	"fmt"
	cacherpccli "seckill/cacheservice/rpc"
	dbrpccli "seckill/dbservice/rpc"
	"seckill/seetings"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type RegisterCenter struct {
	cli      *clientv3.Client
	addr     map[string][]string
	mu       sync.Mutex
	dbCli    map[string]*dbrpccli.ServClient
	cacheCli map[string]*cacherpccli.CacheServCli
	dbidx    int
	cacheidx int
}

var center *RegisterCenter

var once = new(sync.Once)

func GetRegisterCenter() *RegisterCenter {
	once.Do(func() {
		setting, err := seetings.GetSetting()
		if err != nil {
			panic(err)
		}
		addr := fmt.Sprintf("%s:%d", setting.RegisterCenter.Host, setting.RegisterCenter.Port)
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{addr},
			DialTimeout: time.Duration(setting.RegisterCenter.Timeout) * time.Second,
		})
		if err != nil {
			panic(err)
		}
		center = &RegisterCenter{
			cli:      cli,
			addr:     make(map[string][]string),
			mu:       sync.Mutex{},
			dbCli:    make(map[string]*dbrpccli.ServClient),
			cacheCli: make(map[string]*cacherpccli.CacheServCli),
		}
	})
	return center
}

func (rc *RegisterCenter) Register(key string, value string) error {
	fmt.Printf("注册key%s\n", key)
	ctx := context.Background()
	lcr, err := rc.cli.Lease.Grant(ctx, 10)
	if err != nil {
		fmt.Println(err)
		return err
	}
	respChan, err := rc.cli.Lease.KeepAlive(ctx, lcr.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	go func() {
		for {
			select {
			case resp := <-respChan:
				if resp == nil {
					fmt.Printf("service %s closed\n", key)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	kv := clientv3.NewKV(rc.cli)
	_, err = kv.Put(ctx, key, value, clientv3.WithLease(clientv3.LeaseID(lcr.ID)))
	return err
}

func (rc *RegisterCenter) Discovery(ch chan error) {
	kv := clientv3.NewKV(rc.cli)
	ctx := context.Background()
	for {
		dbres, err := kv.Get(ctx, "/bwsk/dbservice/", clientv3.WithPrefix())
		if err != nil {
			ch <- err
			return
		}
		cacheres, err := kv.Get(ctx, "/bwsk/cacheservice/", clientv3.WithPrefix())
		if err != nil {
			ch <- err
			return

		}
		rc.mu.Lock()
		olddb := rc.addr["dbservice"]
		oldcache := rc.addr["cacheservice"]
		rc.addr["dbservice"] = make([]string, 0)
		rc.addr["cacheservice"] = make([]string, 0)
		for _, v := range dbres.Kvs {
			port := string(v.Value)
			if _, ok := rc.dbCli[port]; !ok {
				cli, err := dbrpccli.NewDbServRpcCli(port)
				if err == nil {
					rc.dbCli[port] = cli
					rc.addr["dbservice"] = append(rc.addr["dbservice"], port)
				}
			} else {
				rc.addr["dbservice"] = append(rc.addr["dbservice"], port)
			}
		}
		for _, v := range cacheres.Kvs {
			port := string(v.Value)
			if _, ok := rc.cacheCli[port]; !ok {
				cli, err := cacherpccli.NewCacheServClient(port)
				if err == nil {
					rc.cacheCli[port] = cli
					rc.addr["cacheservice"] = append(rc.addr["cacheservice"], port)
				}
			} else {
				rc.addr["cacheservice"] = append(rc.addr["cacheservice"], port)
			}
		}
		for _, old := range olddb {
			exist := false
			for _, new := range rc.addr["dbservice"] {
				if old == new {
					exist = true
					break
				}
			}
			if !exist {
				delete(rc.dbCli, old)
			}
		}
		for _, old := range oldcache {
			exist := false
			for _, new := range rc.addr["cacheservice"] {
				if old == new {
					exist = true
					break
				}
			}
			if !exist {
				delete(rc.cacheCli, old)
			}
		}
		rc.mu.Unlock()
		c := time.After(5 * time.Second)
		<-c
	}
}

func (rc *RegisterCenter) GetDbClient() (*dbrpccli.ServClient, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if len(rc.addr["dbservice"]) == 0 {
		return nil, errors.New("no db service")
	}
	idx := rc.dbidx % len(rc.addr["dbservice"])
	rc.dbidx++
	cli := rc.dbCli[rc.addr["dbservice"][idx]]
	return cli, nil
}

func (rc *RegisterCenter) GetCacheClient() (*cacherpccli.CacheServCli, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if len(rc.addr["cacheservice"]) == 0 {
		return nil, errors.New("no cache service")
	}
	idx := rc.cacheidx % len(rc.addr["cacheservice"])
	rc.cacheidx++
	cli := rc.cacheCli[rc.addr["cacheservice"][idx]]
	return cli, nil
}
