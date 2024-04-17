package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
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
	service.NewUser().Online(user_id)

	//	用户断开销毁
	defer func() {
		conn.Close()
		// 连接断开时也要销毁连接集里的对象
		service.NewUser().ConnDisconnect(user_id, conn)
	}()

	for {
		var msg model.ConnMsg
		//	ReadJSON 获取值的方式类似于gin的 ctx.ShouldBind() 通过结构体的json映射值
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
		// do something.....

		msg.FromUserID = user_id
		//	发送回信息
		err = conn.WriteJSON(msg)
		if err != nil {
			fmt.Println("用户断开:", err.Error())
			return
		}
		if err = valMsg(msg); err != nil {
			// todo 返回错误提示,但连接不该断
			continue
		}
		// 将消息写入通道
		service.NewChatRoomThread().SendMsg(msg)

	}

}

// 验证数据 例如用户是否有加入聊天室
func valMsg(msg model.ConnMsg) error {
	// do something...
	return nil
}
