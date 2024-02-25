package main

import (
	"context"
	"fmt"
	"hello_go/mxshop/goods_srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var goodsClient proto.GoodsClient
var conn *grpc.ClientConn

func TestGetCategoryList() {
	rsp, err := goodsClient.GetAllCategorysList(context.Background(), &empty.Empty{})

	if err != nil {
		panic(err)
	}
	fmt.Println("Totoal brand number:", rsp.Total)
	fmt.Println(rsp.JsonData)
}

func TestGetSubCategoryList(){
	rsp, err := goodsClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 130358,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(rsp.SubCategorys)
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
	// TestGetCategoryList()
	TestGetSubCategoryList()
	conn.Close()
}
