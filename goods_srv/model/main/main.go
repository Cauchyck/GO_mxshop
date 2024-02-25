package main

import (
	"context"
	"crypto/sha512"
	"fmt"
	"hello_go/mxshop/goods_srv/model"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func genMd5(code string) string {
	// Md5 := md5.New()
	// io.WriteString(Md5, code)
	// return hex.EncodeToString(Md5.Sum(nil))
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("generic password", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newPassword)
	passwordInfo := strings.Split(newPassword, "$")
	cheak := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	fmt.Println(cheak)
	return newPassword
}

func main() {
	// dsn := "root:@tcp(127.0.0.1:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold:             time.Second, // Slow SQL threshold
	// 		LogLevel:                  logger.Info, // Log level
	// 		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
	// 		ParameterizedQueries:      true,        // Don't include params in the SQL log
	// 		Colorful:                  true,        // Disable color
	// 	},
	// )

	// // Globally mode
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	NamingStrategy: schema.NamingStrategy{
	// 		SingularTable: true,
	// 	},
	// 	Logger: newLogger,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// // schema
	// _ = db.AutoMigrate(&model.Category{}, &model.Brands{}, &model.GoodsCategoryBrand{}, &model.Banner{}, &model.Goods{})

	Mysql2Es()
}

func Mysql2Es() {

	dsn := "root:@tcp(127.0.0.1:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	// Globally mode
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "mxshop", log.LstdFlags)

	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"), elastic.SetSniff(false), elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}

	var goods []model.Goods
	db.Find(&goods)

	for _, good := range goods {
		esModel := model.EsGoods{
			ID: good.ID,
			Category: good.CategoryID,
			BrandsID: good.BrandsID,
			OnSale: good.OnSale,
			ShipFree: good.ShipFree,
			IsNew: good.IsNew,
			IsHot: good.IsHot,
			Name: good.Name,
			ClickNum: good.ClickNum,
			SoldNum: good.SoldNum,
			FavNum: good.FavNum,
			MarketPrice: good.MarketPrice,
			GoodsBrief: good.GoodsBrief,
			ShopPrice: good.ShopPrice,
		}

		_, err = client.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(good.ID))).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
