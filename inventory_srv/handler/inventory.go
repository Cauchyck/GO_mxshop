package handler

import (
	"context"
	"fmt"
	"hello_go/mxshop/inventory_srv/global"
	"hello_go/mxshop/inventory_srv/model"
	"hello_go/mxshop/inventory_srv/proto"
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (i *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)

	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "no Inventoty message")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

var m sync.Mutex

func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 数据库基本应用场景：数据库事务
	// 并发场景
	tx := global.DB.Begin()

	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	for _, goodInfo := range req.GoodsInfo {
		// m.Lock()
		var inv model.Inventory

		// Todo 悲观锁
		// if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		// 	tx.Rollback()
		// 	return nil, status.Errorf(codes.InvalidArgument, "no inventory message")
		// }
		// for {
		// 	if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		// 		tx.Rollback()
		// 		return nil, status.Errorf(codes.InvalidArgument, "no inventory message")
		// 	}

		// Todo 乐观锁
		// 	if inv.Stocks < goodInfo.Num {
		// 		tx.Rollback()
		// 		return nil, status.Errorf(codes.ResourceExhausted, "inventory not enough")
		// 	}
		//
		// 	inv.Stocks -= goodInfo.Num
		// 	// tx.Save(&inv)
		// 	result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version = ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1})
		// 	if result.RowsAffected == 0 {
		// 		zap.S().Info("Inventory update faild")
		// 	} else {
		// 		break
		// 	}
		// }

		// Todo redis分布式锁

		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "Internal error")
		}
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "no inventory message")
		}

		if inv.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "inventory not enough")
		}

		inv.Stocks -= goodInfo.Num 
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "release redis lock error")
		}

	}

	tx.Commit()
	// m.Unlock()
	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		m.Lock()
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "no inventory message")
		}

		inv.Stocks += goodInfo.Num
		tx.Save(&inv)

	}

	tx.Commit()
	m.Unlock()

	return &emptypb.Empty{}, nil
}
