package main

import (
	"seckill/registercenter/registerservice"
)

func main() {
	ch := make(chan error, 0)
	registerservice.GetRegisterCenter().Discovery(ch)
	err := <-ch
	panic(err)
}
