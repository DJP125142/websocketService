package service

import (
	"go.uber.org/zap"
	"sync"
	"time"
	"websocketService/global"
	"websocketService/model"
)

type chatRoomThread struct {
	msgChannel chan model.ConnMsg
	lock       sync.Mutex
	once       sync.Once
}

var crt = new(chatRoomThread)

func NewChatRoomThread() *chatRoomThread {
	crt.once.Do(func() {
		crt.msgChannel = make(chan model.ConnMsg, 30) // 通道大小30，可以存30条消息
		crt.lock = sync.Mutex{}                       // todo 写法多余，后续看看能否干掉
	})
	return crt
}

// 启动聊天通道
func (chat *chatRoomThread) Start() {
	// 无限循环的获取消息通道内消息
	for {
		select {
		case msg := <-chat.msgChannel:
			// 表明下发送方user_id
			msg.Msg.Data["from_user_id"] = msg.FromUserID
			if msg.Msg.ChatMsgType == 1 {
				// 群聊消息
				// todo 消息入库
				NewRoom().SendMsgToRoom(msg.Msg.Data["room_id"].(int), msg.Msg)
			} else if msg.Msg.ChatMsgType == 2 {
				// 私聊消息
				NewUser().SendMsgToUid(msg.Msg.Data["to_user_id"].(int), msg.Msg)
			}
		}
	}
}

// 向通道内发送消息
func (chat *chatRoomThread) SendMsg(msg model.ConnMsg) {
	global.Lg.Info("SendMsg", zap.Any("chat", msg))
	//加上时间戳
	msg.Msg.Data["created_at"] = time.Now().Format(time.RFC3339) // 使用 RFC3339 格式化时间字符串
	chat.msgChannel <- msg
}
