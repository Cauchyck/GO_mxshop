package handler

import (
	"context"
	"hello_go/mxshop/userop_srv/global"
	"hello_go/mxshop/userop_srv/model"
	"hello_go/mxshop/userop_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AddressServer struct {
	proto.UnimplementedAddressServer
}

func (*UserOpServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var addresses []model.Address
	var rsp proto.AddressListResponse
	var addressResponse []*proto.AddressResponse

	if result := global.DB.Where(&model.Address{User: req.UserId}).Find(&addresses); result.RowsAffected != 0 {
		rsp.Total = int32(result.RowsAffected)
	}

	for _, address := range addresses {
		addressResponse = append(addressResponse, &proto.AddressResponse{
			Id:           address.ID,
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SingerName,
			SignerMobile: address.SingerMobile,
		})
	}
	rsp.Data = addressResponse

	return &rsp, nil

}

func (*UserOpServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	address := model.Address{
		User:         req.UserId,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SingerName:   req.SignerName,
		SingerMobile: req.SignerMobile,
	}
	global.DB.Save(&address)

	return &proto.AddressResponse{Id: address.ID}, nil

}

func ((*UserOpServer)) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	
	var address model.Address
	
	if result := global.DB.Where("id=? and uesr=?", req.Id, req.UserId).First(&address); result.RowsAffected == 0{
		return nil, status.Errorf(codes.NotFound, "地址不存在")
	}

	if req.Province != ""{
		address.Province =req.Province
	}
	if req.City != ""{
		address.City =req.City
	}
	if req.District  != "" {
		address.District = req.District
	}
	if req.Address != ""{
		address.Address = req.Address
	}
	if req.SignerName != ""{
		address.SingerName = req.SignerName
	}
	if req.SignerMobile != ""{
		address.SingerMobile = req.SignerMobile
	}

	global.DB.Save(&address)

	return &emptypb.Empty{}, nil
}


func (*UserOpServer) DeletcAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	
	if result := global.DB.Where("id=? and user=?", req.Id, req.UserId).Delete(&model.Address{}); result.RowsAffected == 0{
		return nil, status.Errorf(codes.NotFound, "收货地址不存在")
	}
	return &emptypb.Empty{}, nil
}
