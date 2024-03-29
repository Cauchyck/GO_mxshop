package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"hello_go/mxshop/goods_srv/global"
	"hello_go/mxshop/goods_srv/model"
	"hello_go/mxshop/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 商品分类
func (s *GoodsServer)GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)

	for _, cacategory := range categorys {
		fmt.Println(cacategory.Name)
	}
	b, _ := json.Marshal(&categorys)

	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

// 获取⼦分类
func (s *GoodsServer)GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error){
	categoryListResponse := proto.SubCategoryListResponse{}
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "category has not exist")
	}

	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id: category.ID,
		Name: category.Name,
		Level: category.Level,
		IsTab: category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}

	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	preloads := "SubCategory"
	if category.Level == 1{
		preloads = "SubCategory.SubCategory"
	}
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload(preloads).Find(&subCategorys)

	for _, subCategory := range subCategorys{
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id: subCategory.ID,
			Name: subCategory.Name,
			Level: subCategory.Level,
			IsTab: subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}

	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse, nil
}
func (s *GoodsServer)CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error){
	category := model.Category{
		Name: req.Name,
		Level: req.Level,
		IsTab: req.IsTab,
	}
	if req.Level != 1{
		category.ParentCategoryID = req.ParentCategory
	}
	global.DB.Save(&category)

	return &proto.CategoryInfoResponse{Id: category.ID}, nil
}
func (s *GoodsServer)DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error){

	if result := global.DB.Find(&model.Category{}, req.Id); result.RowsAffected == 0{
		return nil, status.Errorf(codes.NotFound, "goods has not exist")
	}

	return &emptypb.Empty{}, nil
}

func (s *GoodsServer)UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error){
	var catgory model.Category

	if result := global.DB.Find(&catgory, req.Id); result.RowsAffected == 0{
		return nil, status.Errorf(codes.NotFound, "goods has not exist")
	}

	if req.Name != ""{
		catgory.Name = req.Name
	}
	if req.ParentCategory != 0 {
		catgory.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0{
		catgory.Level = req.Level
	}
	if req.IsTab{
		catgory.IsTab = req.IsTab
	}
	global.DB.Save(&catgory)

	return &emptypb.Empty{}, nil
}
