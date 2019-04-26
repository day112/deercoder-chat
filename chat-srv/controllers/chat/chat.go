package chat

import (
	"context"
	"deercoder-chat/chat-srv/models/chat"
	"deercoder-chat/chat-srv/proto"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Streamer struct{}
type ChatService struct{}

// Server side stream
func (e *Streamer) ServerStream(ctx context.Context, req *proto.Request, stream proto.Streamer_ServerStreamStream) error {
	log.Printf("[Chat-srv]: Got msg %v", req.Message)
	if err := stream.Send(&proto.Response{Message: req.Message}); err != nil {
		return err
	}
	return nil
}

// Bidirectional stream
func (e *Streamer) Stream(ctx context.Context, stream proto.Streamer_StreamStream) error {
	for {
		// Read from stream
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Printf("[Chat-srv]:Got msg %v", req.Message)
		if err := stream.Send(&proto.Response{Message: req.Message}); err != nil {
			return err
		}
	}
}

//=========chat service===================================================
//========================================================================
//========================================================================

// 创建群聊
func (c *ChatService) DistributeGroup(ctx context.Context, req *proto.UidS, rsp *proto.Response) error {
	uids := req.Uids
	gid, _ := chat.DistributeGroup(uids)
	if gid == "" {
		return errors.New("群聊创建失败")
	}
	rsp.Message.GroupId = gid
	return nil
}

// 拉取群聊所有消息
func (c *ChatService) GetAllGroupMsg(ctx context.Context, req *proto.Request, rsp *proto.ArrayMessage) error {
	group_id := req.Message.GroupId

	msg, err := chat.GetAllGroupMsg(group_id)
	if err != nil {
		return err
	}
	rsp.Message = msg
	return nil
}

// 拉取离线信息
func (c *ChatService) GetGroupLastMsg(ctx context.Context, req *proto.Request, rsp *proto.ArrayMessage) error {
	group_id := req.Message.GroupId
	uid := req.Message.FromUid

	msg, err := chat.GetGroupLastMsg(group_id, uid)
	if err != nil {
		return err
	}
	rsp.Message = msg
	return nil
}

// 已读离线信息
func (c *ChatService) ReadGroupLastMsg(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	group_id := req.Message.GroupId
	uid := req.Message.FromUid

	_, err := chat.ReadGroupLastMsg(group_id, uid)

	return err
}

// 获取用户好友列表
func (c * ChatService) GetUserList(ctx context.Context, req *proto.ChatUser, rsp *proto.UserList) (err error) {

	//users := []*proto.ChatUser
	rsp.UserList, err = chat.GetUserList(req.Id)
	return err
}

// 获取群聊中用户列表
func (c * ChatService) GetGroupUser(ctx context.Context, req *proto.GroupUser, rsp *proto.GUserResponse) error {
	return chat.GetGroupUser(req.GroupId, rsp.GroupUser)
}

// 群发消息
func MassMessage(u *gin.Context) {

	group_ids := u.PostForm("group_ids")
	send_uids := u.PostForm("send_uids")
	from_uid := u.PostForm("from_uid")
	content := u.PostForm("content")
	ss := chat.MassMessage(group_ids, send_uids, from_uid, content)
	u.JSON(http.StatusOK, ss)
}
