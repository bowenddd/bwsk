package main

import (
	"fmt"
	"seckill/dbservice/rpc"
)

func main() {
	rpcServ := rpc.GetRpcServServer()
	err := rpcServ.StartRpcServServer()
	if err != nil {
		fmt.Println("start rpc serv error")
	}
}
