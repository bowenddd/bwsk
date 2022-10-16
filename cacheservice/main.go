package main

import (
	"seckill/cacheservice/rpc"
)

func main() {
    rpcServ := rpc.GetCacheRpcService()
    rpcServ.StartCacheRpcServService()
}
