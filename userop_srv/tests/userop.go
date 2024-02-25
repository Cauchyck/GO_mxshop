package main

import (
	"context"
	"fmt"
	"hello_go/mxshop/userop_srv/proto"

	"google.golang.org/grpc"
)

var userFavClient proto.UserFavClient
var messageClient proto.MessageClient
var addressClient proto.AddressClient

var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:57934", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	fmt.Println("conn to userop_src")
	userFavClient = proto.NewUserFavClient(conn)
	messageClient = proto.NewMessageClient(conn)
	addressClient = proto.NewAddressClient(conn)
}

func TestAddressList(){
	rsp, err := addressClient.GetAddressList(context.Background(), &proto.AddressRequest{
		UserId: 1,
	})
	if err != nil{
		panic(err)
	}

	fmt.Println(rsp)
}


func TesMessageList(){
	rsp, err := messageClient.GetMessageList(context.Background(), &proto.MessageRequest{
		UserId: 1,
	})
	if err != nil{
		panic(err)
	}

	fmt.Println(rsp)
}

func TesUserfavList(){
	rsp, err := userFavClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: 1,
	})
	if err != nil{
		panic(err)
	}

	fmt.Println(rsp)
}

func main() {
	Init()
	TestAddressList()
	TesMessageList()
	TesUserfavList()
	conn.Close()
}
