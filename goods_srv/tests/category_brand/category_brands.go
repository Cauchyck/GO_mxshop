package main

import (
	"context"
	"fmt"
	"hello_go/mxshop/goods_srv/proto"

	"google.golang.org/grpc"
)

var goodsClient proto.GoodsClient
var conn *grpc.ClientConn

func TestCategoryBrandList() {
	rsp, err := goodsClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})

	if err != nil {
		panic(err)
	}
	fmt.Println("Totoal brand number:", rsp.Total)
	fmt.Println(rsp.Data)
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:58923", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	fmt.Println("conn to goods_srv")
	goodsClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()

	conn.Close()
}
