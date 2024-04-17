package model

// 消息通道
type ConnMsg struct {
	Msg        ChatMsg `json:"msg,omitempty"`
	FromUserID int     `json:"from_user_id,omitempty"`
}

// 消息内容
// ChatMsgType = 1 群聊信息 ChatMsgType = 2 一对一信息 ...
type ChatMsg struct {
	ChatMsgType int                    `json:"chat_msg_type,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}
