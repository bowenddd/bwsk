package main

import (
	"fmt"
	"seckill/dbservice/rpc/serv"
	"sync"
)

func main() {
	ports := []string{":53277", ":53278", ":53279"}
	wg := new(sync.WaitGroup)
	for i := 0; i < len(ports); i++ {
		wg.Add(1)
		go func(port string) {
			defer wg.Done()
			serv := serv.NewRpcServServer()
			err := serv.StartRpcServServer(port)
			if err != nil {
				fmt.Println(err)
			}
		}(ports[i])
	}
	wg.Wait()
}
