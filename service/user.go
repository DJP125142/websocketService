package service

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"websocketService/model"
)

func init() {
	NewUser()
}

// 定义一个map来存储每个用户的id和对应的ws连接
type userConn struct {
	userConn map[int]*websocket.Conn
	lock     sync.Mutex // 互斥锁，保障map的并发安全
	once     sync.Once  // 保障初始化只操作一次
}

// 声明一个用户对象
var user = new(userConn)

// 单例模式
func NewUser() *userConn {
	// 只执行一次初始化
	user.once.Do(func() {
		user.userConn = make(map[int]*websocket.Conn)
		user.userConn[-1] = nil
		user.lock = sync.Mutex{}
	})
	return user
}

// 用户发起websocket连接
// join_type 加入模式
//
//	1 正常加入 占线无法加入
//	2 强制加入 即踢下线前者
func (user *userConn) ConnConnect(user_id, join_type int, conn *websocket.Conn) (int, error) {
	user.lock.Lock()
	defer user.lock.Unlock()

	// 判断加入方式
	if join_type == 1 {
		// 判断用户是否已经在线
		if _, ok := user.userConn[user_id]; ok {
			return 1, errors.New("该账号已被登录")
		}
	} else if join_type == 2 {
		// 如果用户已经存在map内，进行销毁挤出
		if conn2, ok := user.userConn[user_id]; ok {
			err := conn2.Close()
			if err != nil {
				fmt.Println(err)
			}
			delete(user.userConn, user_id)
		}
		// 重新加入
		user.userConn[user_id] = conn
	}
	return -1, nil
}

// 断开连接
func (user *userConn) ConnDisconnect(user_id int, conn *websocket.Conn) error {
	user.lock.Lock()
	defer user.lock.Unlock()

	if conn2, ok := user.userConn[user_id]; ok {
		if conn == conn2 {
			delete(user.userConn, user_id)
		}
	} else {
		// 不存在的连接申请断开
		// todo add error log
	}
	return nil
}

// 对单个用户发送消息
func (user *userConn) SendMsgToUid(user_id int, msg interface{}) error {
	var err error
	// 获取目标用户的消息通道
	if conn, ok := user.userConn[user_id]; ok {
		err = conn.WriteJSON(msg) // 给目标user_id通道里写入消息
	} else {
		err = errors.New("该用户不在线")
	}
	return err
}

// 对多个用户发送消息
func (user *userConn) SendMsgToUidList(user_ids []int, msg interface{}) (res_user_ids []int, res_errs []error) {
	for _, user_id := range user_ids {
		content := msg.(model.ChatMsg)
		if content.ChatMsgType == 1 {
			// 群聊中遍历到自己时，跳过发送
			if user_id == content.Data["from_user_id"].(int) {
				continue
			}
		}
		// 获取每个目标用户的消息通道
		if conn, ok := user.userConn[user_id]; ok {
			err := conn.WriteJSON(msg)
			if err != nil {
				res_user_ids = append(res_user_ids, user_id)
				res_errs = append(res_errs, err)
			}
		} else {
			res_user_ids = append(res_user_ids, user_id)
			res_errs = append(res_errs, errors.New("该用户不在线"))
		}
	}
	return
}

// 用户上线通知
// 通知大厅
func (user *userConn) Online(user_id int) {
	// 构建一条系统的上线通知消息
	var content model.ChatMsg
	content.ChatMsgType = 1
	content.Data = map[string]interface{}{
		"room_id": 1,
		"content": "新用户已上线",
	}

	var msg model.ConnMsg
	msg.FromUserID = user_id
	msg.Msg = content

	NewChatRoomThread().SendMsg(msg)

}
