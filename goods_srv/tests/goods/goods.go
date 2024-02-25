package main

import (
	"context"
	"fmt"
	"hello_go/mxshop/goods_srv/proto"

	"google.golang.org/grpc"
)

var goodsClient proto.GoodsClient
var conn *grpc.ClientConn

func TestGoodsList() {
	rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		PriceMin: 90,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println("Totoal brand number:", rsp.Total)


	for _, good := range rsp.Data{
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func TestBatchGetGoods() {
	rsp, err := goodsClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: []int32{421, 422, 423},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println("Totoal brand number:", rsp.Total)


	for _, good := range rsp.Data{
		fmt.Println(good.Name, good.ShopPrice)
	}
}
func TestGetGoodsDetail() {
	rsp, err := goodsClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: 421,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Name)
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:54617", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	fmt.Println("conn to goods_srv")
	goodsClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()

	// TestGoodsList()
	TestBatchGetGoods()
	TestGetGoodsDetail()
	conn.Close()
}
