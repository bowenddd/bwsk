package main

import (
	"fmt"
	"seckill/cacheservice/rpc"
)

func main() {
	csc, err := rpc.NewCacheServClient()

	if err != nil {
		panic(err)
	}
	i, err := csc.GetStock(165)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(i)
}
