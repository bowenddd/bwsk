package main

import (
	"fmt"
	"seckill/dbservice/service"
)

func main() {
	productServ := service.GetProductServ()
	err := productServ.DeleteProduct("test1")
	fmt.Println(err)

}
