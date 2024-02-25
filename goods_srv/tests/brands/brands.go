package main

import (
	"context"
	"fmt"
	"hello_go/mxshop/goods_srv/proto"

	"google.golang.org/grpc"
)

var goodsClient proto.GoodsClient
var conn *grpc.ClientConn

func TestGetBrandList() {
	rsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{})

	if err != nil {
		panic(err)
	}
	fmt.Println("Totoal brand number:", rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
		if err != nil {
			panic(err)
		}
	}
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:62489", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	fmt.Println("conn to goods_srv")
	goodsClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()
	TestGetBrandList()
	conn.Close()
}
