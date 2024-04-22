package service

import (
	"go.uber.org/zap"
	"sync"
	"websocketService/global"
)

// 定义一个map来存储群内的用户id
type rooms struct {
	members map[int]map[int]struct{}
	lock    sync.Mutex
	once    sync.Once
}

var room = new(rooms)

// 初始化房间
func NewRoom() *rooms {
	room.once.Do(func() {
		room.members = make(map[int]map[int]struct{})
		room.lock = sync.Mutex{}
	})
	return room
}

// 生成一个新的房间id
func (room *rooms) CreateRoomId() int {
	room.lock.Lock()
	defer room.lock.Unlock()
	// 找到当前有活跃的最大房间号
	curMaxRoomId := 1
	for key := range room.members {
		if key > curMaxRoomId {
			curMaxRoomId = key
		}
	}
	// 生成一个max+1的新房间号
	newMaxRoomId := curMaxRoomId + 1
	return newMaxRoomId
}

// 获取房间内在线人数
func (room *rooms) GetRoomOnlineUserCount(room_id int) int {
	room.lock.Lock()
	defer room.lock.Unlock()
	return len(room.members[room_id])
}

// 发送群消息
func (room *rooms) SendMsgToRoom(room_id int, msg interface{}) {
	user := NewUser()
	room.lock.Lock()
	defer room.lock.Unlock()

	global.Lg.Info("SendMsgToRoom", zap.Any("room_id", room_id))
	global.Lg.Info("SendMsgToRoom", zap.Any("members", room.members))

	// 获取房间内所有人user_id
	var user_ids []int
	for key, _ := range room.members[room_id] {
		user_ids = append(user_ids, key)
	}
	global.Lg.Info("SendMsgToRoom", zap.Any("user_ids", user_ids))
	// 批量发送消息
	user.SendMsgToUidList(user_ids, msg)
}

// 用户上线，同时加入多个聊天室
func (room *rooms) UserJoinRooms(room_ids []int, user_id int) {
	room.lock.Lock()
	defer room.lock.Unlock()
	for _, room_id := range room_ids {
		if v, ok := room.members[room_id]; !ok {
			// map中没有说明房间不存在，则新建房间
			room.members[room_id] = make(map[int]struct{})
			room.members[room_id][user_id] = struct{}{}
		} else {
			v[user_id] = struct{}{}
		}
	}
	return
}

// 用户下线/退群 退出聊天室链接集合
func (room *rooms) UserQuitRooms(room_ids []int, user_id int) {
	room.lock.Lock()
	defer room.lock.Unlock()
	for _, room_id := range room_ids {
		if v, ok := room.members[room_id]; ok {
			delete(v, user_id)
			//	房间没人就销毁
			if len(room.members[room_id]) <= 0 {
				delete(room.members, room_id)
			}
		}
	}
	return
}

// 用户上线/入群 加入聊天室连接集合
func (room *rooms) UserJoinRoom(room_id, user_id int) {
	room.lock.Lock()
	defer room.lock.Unlock()
	if v, ok := room.members[room_id]; !ok {
		//	房间不存在就创建
		room.members[room_id] = make(map[int]struct{})
		room.members[room_id][user_id] = struct{}{}
	} else {
		// 存在则加入
		v[user_id] = struct{}{}
	}
	return
}

// 用户下线/退群 退出聊天室链接集合
func (room *rooms) UserQuitRoom(room_id, user_id int) {
	room.lock.Lock()
	defer room.lock.Unlock()
	if v, ok := room.members[room_id]; ok {
		delete(v, user_id)
		//	房间没人就销毁
		if len(room.members[room_id]) <= 0 {
			delete(room.members, room_id)
		}
	}
	return
}

// 用户上线/入群 加入聊天室连接集合
func (room *rooms) UsersJoinRoom(room_id int, user_ids []int) {
	room.lock.Lock()
	defer room.lock.Unlock()
	if v, ok := room.members[room_id]; !ok {
		//	房间不存在就创建
		room.members[room_id] = make(map[int]struct{})
		for _, user_id := range user_ids {
			room.members[room_id][user_id] = struct{}{}
		}
	} else {
		// 存在则加入
		for _, user_id := range user_ids {
			v[user_id] = struct{}{}
		}
	}
	return
}
