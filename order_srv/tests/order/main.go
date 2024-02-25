package main

import (
	"context"
	"fmt"
	"hello_go/mxshop/order_srv/proto"

	"google.golang.org/grpc"
)

var orderClient proto.OrderClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:54597", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	orderClient = proto.NewOrderClient(conn)
}

func TestCreateCartItem(userId, nums, goodsId int32) {
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		Nums:    nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}

func TestCartItemList(userId int32) {
	rsp, err := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id:  userId,
	})
	if err != nil {
		panic(err)
	}
	for _, item := range rsp.Data{
		fmt.Println(item.Id, item.GoodsId, item.Nums)
	}
}

func TestUpdateCartItem(id int32) {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:  id,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}

}

func TestCreateOrder(){
	_, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId: 1,
		Address: "shanghai",
		Name: "bobby",
		Mobile: "18734934893",
		Post: "!!!",
	})

	if err != nil{
		panic(err)
	}

}

func TestOrderDetail(orderId int32){
	rsp, err := orderClient.OrderDetail(context.Background(), &proto.OrderRequest{
		Id: orderId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.OrderInfo.OrderSn)
	for _, good := range rsp.Goods {
		fmt.Println(good.GoodsName)
	}
}

func TestOrderList(){
	rsp, err := orderClient.OrderList(context.Background(), &proto.OrderFilterRequest{
		UserId: 0,
	})
	if err != nil {
		panic(err)
	}
	for _, order := range rsp.Data {
		fmt.Println(order.OrderSn)
	}
}

func main() {
	Init()
	// TestCreateCartItem(1, 1, 422)
	// TestCartItemList(1)
	// TestUpdateCartItem(1)
	// TestCreateOrder()
	// TestOrderDetail(2)
	TestOrderList()
	conn.Close()
}
