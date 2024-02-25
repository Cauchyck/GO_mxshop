package handler

import (
	"context"
	"hello_go/mxshop/userop_srv/global"
	"hello_go/mxshop/userop_srv/model"
	"hello_go/mxshop/userop_srv/proto"
)

type MessageServer struct {
	proto.UnimplementedMessageServer
}

func (*UserOpServer) GetMessageList(ctx context.Context, req *proto.MessageRequest) (*proto.MessageListResponse, error) {
	var messages []model.LeavingMessage
	var rsp proto.MessageListResponse
	var messageResponse []*proto.MessageResponse

	if result := global.DB.Where(&model.LeavingMessage{User: req.UserId}).Find(&messages); result.RowsAffected != 0 {
		rsp.Total = int32(result.RowsAffected)
	}

	for _, message := range messages {
		messageResponse = append(messageResponse, &proto.MessageResponse{
			Id:          message.ID,
			UserId:      message.User,
			MessageType: message.MessageType,
			Subject:     message.Subject,
			Message:     message.Message,
			File:        message.File,
		})
	}
	rsp.Data = messageResponse

	return &rsp, nil
}
func (*UserOpServer) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	message := model.LeavingMessage{
		User:        req.UserId,
		MessageType: req.MessageType,
		Subject:     req.Subject,
		Message:     req.Message,
		File:        req.File,
	}

	global.DB.Save(&message)

	return &proto.MessageResponse{Id: message.ID}, nil
}
