package main

import (
	"seckill/cacheservice/rpc/serv"
	"sync"
)

func main() {
	ports := []string{":46288", ":46289", ":46290"}
	wg := new(sync.WaitGroup)
	for i := 0; i < len(ports); i++ {
		wg.Add(1)
		go func(port string) {
			defer wg.Done()
			serv := serv.NewCacheRpcService()
			serv.StartCacheRpcServService(port)
		}(ports[i])
	}
	wg.Wait()
}
