package main

import (
	"context"
	"fmt"
	"hello_go/mxshop/user_srv/proto"

	"google.golang.org/grpc"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:8888", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	fmt.Println("conn to 8888")
	userClient = proto.NewUserClient(conn)
}
func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 2,
	})

	if err != nil {
		panic(err)
	}

	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord)
		checkRsp, err := userClient.CheckPassWord(context.Background(), &proto.PassWordCheckInfo{
			PassWord:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}

func TestCreateUser() {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("noppy%d", i),
			Mobile:   fmt.Sprintf("1323432343%d", i),
			PassWord: "admin123",
		})

		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}

}

func main() {
	Init()
	// TestGetUserList()
	TestCreateUser()
	conn.Close()
}
