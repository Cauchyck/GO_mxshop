package main

import (
	"crypto/sha512"
	"fmt"
	"hello_go/mxshop/user_srv/model"
	"log"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
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
	dsn := "root:@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

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
	// schema
	// _ = db.AutoMigrate(&model.User{})

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("generic password", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	for i := 0; i < 10; i++ {
		user := model.User{
			NickName: fmt.Sprintf("bobby%d", i),
			Mobile:   fmt.Sprintf("1877887787%d", i),
			Password: newPassword,
		}
		db.Save(&user)
	}

}
