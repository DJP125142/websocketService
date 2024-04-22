package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
	"websocketService/global"
	"websocketService/model"
	"websocketService/response"
	"websocketService/service"
	"websocketService/utils"
)

// websocket配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

var ModelChatRoom model.ChatRoom

// 用户申请创建socket链接
func CreateConn(c *gin.Context) {
	var (
		conn    *websocket.Conn
		err     error
		user_id int
	)

	// 从token里获取user_id
	token := c.Query("token")
	userInfoData, err := utils.GetUserInfoByToken(token)
	if err != nil {
		response.Err(c, 401, "获取用户信息失败", err.Error())
	}
	userInfo, _ := userInfoData["userInfo"].(map[string]interface{})
	user_id = int(userInfo["userId"].(float64))

	//	判断请求过来的链接是否要升级为websocket
	if websocket.IsWebSocketUpgrade(c.Request) {
		//	将请求升级为 websocket链接
		conn, err = upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
		if err != nil {
			response.Err(c, 500, "创建链接失败", err.Error())
			return
		}
	} else {
		response.Err(c, 500, "创建链接失败", "")
		return
	}

	// 获取用户加入的聊天室id数组
	//room_ids, _ := ModelChatRoom.GetUserRoomIds(user_id)

	// 用户加入大厅
	service.NewRoom().UserJoinRoom(1, user_id)
	// 用户加入连接集,用强制加入方式
	service.NewUser().ConnConnect(user_id, 2, conn)
	// 用户上线通知
	service.NewUser().Online(user_id, userInfo["username"].(string))

	//	用户断开销毁
	defer func() {
		conn.Close()
		// 用户离开房间
		service.NewRoom().UserQuitRoom(1, user_id)
		// 连接断开时也要销毁连接集里的对象
		service.NewUser().ConnDisconnect(user_id, conn)
		// 用户下线通知
		service.NewUser().Offline(user_id, userInfo["username"].(string))
	}()

	for {
		var msg model.ConnMsg
		//	如果读不到值 则堵塞在此处
		err = conn.ReadJSON(&msg)
		if err != nil {
			// 写回错误信息
			err = conn.WriteJSON(map[string]interface{}{"code": 500, "msg": "获取数据失败", "data": ""})
			if err != nil {
				fmt.Println("用户断开", err.Error())
				return
			}
		}

		msg.FromUserID = user_id
		msg.Msg.Data["from_user_id"] = user_id
		if err = valMsg(&msg); err != nil {
			// todo 返回错误提示,但连接不该断
			continue
		}

		//	发送回信息
		//err = conn.WriteJSON(msg.Msg)
		// 将消息写入通道
		service.NewChatRoomThread().SendMsg(msg)

	}

}

// 验证数据 例如用户是否有加入聊天室
func valMsg(msg *model.ConnMsg) error {
	global.Lg.Info("valMsg", zap.Any("msg", msg))
	if msg.Msg.MsgType == 3 {
		user_ids := []int{msg.Msg.Data["from_user_id"].(int), int(msg.Msg.Data["to_user_id"].(float64))}
		room_id := service.NewRoom().CreateRoomId()
		msg.Msg.Data["room_id"] = room_id
		service.NewRoom().UsersJoinRoom(room_id, user_ids)
	}
	global.Lg.Info("valMsg", zap.Any("createRoomId", msg))
	if err := validateAndConvertData(msg.Msg.Data); err != nil {
		return err
	}
	//加上时间戳
	msg.Msg.Data["created_at"] = time.Now().Format(time.RFC3339) // 使用 RFC3339 格式化时间字符串
	return nil
}

// validateAndConvertData 检查并转换 Data 中的 room_id 为正确的 int 类型
func validateAndConvertData(data map[string]interface{}) error {
	if roomId, ok := data["room_id"]; ok {
		switch v := roomId.(type) {
		case float64:
			// JSON 解码可能将数值解释为 float64，需要转换为 int
			data["room_id"] = int(v)
		case int:
			// room_id 已经是 int 类型，无需转换
		default:
			return errors.New("room_id must be a number")
		}
	}
	return nil
}
