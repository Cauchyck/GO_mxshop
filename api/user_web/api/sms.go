package api

import (
	"context"
	"fmt"
	"hello_go/mxshop/api/user_web/forms"
	"hello_go/mxshop/api/user_web/global"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func GenerateSmsCode(witdh int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)

	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {

	zap.S().Info("[forms.PassWordLoginForm{}]")

	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	smsCode := GenerateSmsCode(6)
	fmt.Println("send smsCode: ", smsCode, " to ", sendSmsForm.Mobile)
	
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, 300*time.Second)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "send success",
	})
}
