package message

import (
	"context"
	"hello_go/mxshop/api/userop_web/api"
	"hello_go/mxshop/api/userop_web/forms"
	"hello_go/mxshop/api/userop_web/global"
	"hello_go/mxshop/api/userop_web/models"
	"hello_go/mxshop/api/userop_web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetMessageList(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	request := proto.MessageRequest{}

	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.MessageSrvClient.GetMessageList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	messageList := make([]interface{}, 0)

	for _, item := range rsp.Data {
		messageList = append(messageList, map[string]interface{}{
			"id":      item.Id,
			"user_id": item.UserId,
			"type":    item.MessageType,
			"subject": item.Subject,
			"message": item.Message,
			"file":    item.File,
		})
	}
	reMap["data"] = messageList
	ctx.JSON(http.StatusOK, reMap)

}

func NewMessage(ctx *gin.Context) {
	messageForm := forms.MessageForm{}
	if err := ctx.ShouldBindJSON(&messageForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}
	userId, _ := ctx.Get("userId")

	rsp, err := global.MessageSrvClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})

	if err != nil {
		zap.S().Errorw("新建失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})

}
