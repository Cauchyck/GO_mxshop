package main

import (
	"context"
	"fmt"
	"hello_go/mxshop/inventory_srv/proto"
	"sync"

	"google.golang.org/grpc"
)

var invClient proto.InventoryClient
var conn *grpc.ClientConn

func Init()  {
	var err error
	conn, err = grpc.Dial("127.0.0.1:61278", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	invClient = proto.NewInventoryClient(conn)
}

func TestSetInv(goodsId, num int32) {
	_, err := invClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
		Num:     num,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Set Inventory Number Success")
}

func TestInvDetail(goodsId int32) {
	rsp, err := invClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Num)
}

func TestSell(wg *sync.WaitGroup) {

	defer wg.Done()
	_, err := invClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{
				GoodsId: 421,
				Num:     1,
			},
		},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println("Reset Inventory Number Success")
}

func TestReback() {
	_, err := invClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{
				GoodsId: 421,
				Num:     1,
			},
		},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println("Reback Inventory Number Success")
}
func main() {
	Init()
	// for i := 421; i <= 840; i++{
	// 	TestSetInv(int32(i), 100)
	// }
	// TestSetInv(422, 40)
	// TestInvDetail(421)
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i<20; i++{
		go TestSell(&wg)
	}
	wg.Wait()
	// TestSell()
	// TestReback()
	conn.Close()
}
