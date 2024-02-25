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


func (*UserOpServer) GetFavList(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	var userFavs []model.UserFav
	var rsp proto.UserFavListResponse
	var userFavResponse []*proto.UserFavResponse
	// 查询用户的收藏记录
	// 查询某件商品被哪些用户收藏
	if result := global.DB.Where(&model.UserFav{User: req.UserId, Goods: req.GoodsId}).Find(&userFavs); result.RowsAffected != 0 {
		rsp.Total = int32(result.RowsAffected)
	}

	for _, userFav := range userFavs {
		userFavResponse = append(userFavResponse, &proto.UserFavResponse{
			UserId:  userFav.ID,
			GoodsId: userFav.Goods,
		})
	}
	rsp.Data = userFavResponse

	return &rsp, nil
}
func (*UserOpServer) AddUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {

	global.DB.Save(&model.UserFav{
		User:  req.UserId,
		Goods: req.GoodsId,
	})

	return &emptypb.Empty{}, nil
}
func (*UserOpServer) DelectUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	// Unscoped：物理删除，避免再次收藏时的联合索引冲突
	if result := global.DB.Unscoped().Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.UserFav{}); result.RowsAffected==0{
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}
	return &emptypb.Empty{}, nil
}
func (*UserOpServer) GetUserFavDetail(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	var userFav model.UserFav
	if result := global.DB.Where("goods=? and user=?",req.GoodsId,req.UserId).Find(&userFav); result.RowsAffected == 0{
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}
	return &emptypb.Empty{}, nil
}
