package handler

import (
	"context"
	"hello_go/mxshop/goods_srv/global"
	"hello_go/mxshop/goods_srv/model"
	"hello_go/mxshop/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌分类
func (s *GoodsServer)CategoryBrandList(ctx context.Context, req *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	var categoryBrands []model.GoodsCategoryBrand
	catgoryBrandListResponse := proto.CategoryBrandListResponse{}

	var total int64
	global.DB.Model(&model.GoodsCategoryBrand{}).Count(&total)

	catgoryBrandListResponse.Total = int32(total)

	global.DB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&categoryBrands)

	var categoryResposes []*proto.CategoryBrandResponse
	for _, categoryBrand := range categoryBrands {
		categoryResposes = append(categoryResposes, &proto.CategoryBrandResponse{
			Category: &proto.CategoryInfoResponse{
				Id:             categoryBrand.Category.ID,
				Name:           categoryBrand.Category.Name,
				Level:          categoryBrand.Category.Level,
				IsTab:          categoryBrand.Category.IsTab,
				ParentCategory: categoryBrand.Category.ParentCategoryID,
			},
			Brand: &proto.BrandInfoResponse{
				Id:   categoryBrand.Brands.ID,
				Name: categoryBrand.Brands.Name,
				Logo: categoryBrand.Brands.Logo,
			},
		})
	}
	catgoryBrandListResponse.Data = categoryResposes
	return &catgoryBrandListResponse, nil
}

// 通过category获取brands
func (s *GoodsServer)GetCategoryBrandList(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	brandsListResponse := proto.BrandListResponse{}

	var category model.Category
	if result := global.DB.Find(&category, req.Id).First(&category); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "category is not exist")
	}

	var categoryBrands []model.GoodsCategoryBrand
	if result := global.DB.Preload("Brands").Where(&model.GoodsCategoryBrand{CategoryID: category.ID}).Find(&categoryBrands); result.RowsAffected > 0 {
		brandsListResponse.Total = int32(result.RowsAffected)
	}

	var brandInfoResponses []*proto.BrandInfoResponse
	for _, categoryBrand := range categoryBrands {
		brandInfoResponses = append(brandInfoResponses, &proto.BrandInfoResponse{
			Id:   categoryBrand.Brands.ID,
			Name: categoryBrand.Brands.Name,
			Logo: categoryBrand.Brands.Logo,
		})
	}
	brandsListResponse.Data = brandInfoResponses

	return &brandsListResponse, nil
}
func (s *GoodsServer)CreateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	var categoty model.Category
	if result := global.DB.Find(&categoty, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "goods category has not exist")
	}

	var brand model.Brands
	if result := global.DB.Find(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "goods brands has not exist")
	}
	categoryBrand := model.GoodsCategoryBrand{
		CategoryID: req.CategoryId,
		BrandsID:   req.BrandId,
	}

	global.DB.Save(&categoryBrand)

	return &proto.CategoryBrandResponse{Id: categoryBrand.ID}, nil
}
func (s *GoodsServer)DeleteCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Find(&model.GoodsCategoryBrand{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "Brand category has not exist")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer)UpdateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error){
	var categoryBrand model.GoodsCategoryBrand
	if result := global.DB.First(&categoryBrand, req.Id); result.RowsAffected == 0{
		return nil, status.Errorf(codes.NotFound, "Brand category not exist")
	}
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "brand not exist")
	}
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected ==0 {
		return nil, status.Errorf(codes.InvalidArgument, "category not exist")
	}

	categoryBrand.CategoryID = req.CategoryId
	categoryBrand.BrandsID = req.BrandId

	global.DB.Save(&categoryBrand)

	return &emptypb.Empty{}, nil
}
