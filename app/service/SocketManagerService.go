package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"okc/utils"
	"sync"
	"time"
)

// WsManager
// 客户端管理
var WsManager *WsClientManager

// WsClientManager
// 客户端连接管理
type WsClientManager struct {
	Clients    map[string]*WsClient // 记录在线用户
	Broadcast  chan []byte          // 触发消息广播
	Register   chan *WsClient       // 触发用户登录
	UnRegister chan *WsClient       // 触发用户退出
	Mx         sync.Mutex           // 锁
}

// WsClientBroadcast
// 用户消息广播
func (w *WsClientManager) WsClientBroadcast() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("w *WsClientManager WsClientBroadcast() [ERROR] : %s", err))
			return
		}
		_ = utils.WriteDebugLog(fmt.Sprintf("w *WsClientManager WsClientBroadcast() 用户消息广播协程退出"))
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("w *WsClientManager WsClientBroadcast() 用户消息广播协协程启动"))

	for {
		select {
		case msg, ok := <-w.Broadcast:

			if !ok {
				_ = utils.WriteDebugLog(fmt.Sprintf("w.Broadcast close [ERROR] : %s", ok))
				return
			}
			//log.Println("接收到广播消息:", string(msg))
			// 消息广播
			for uuid, client := range w.Clients {
				// 用户已退出
				if client.IsCloseConn {
					continue
				}
				// 发送用户消息
				//w.Mx.Lock()
				err := client.Socket.WriteMessage(websocket.TextMessage, msg)
				//w.Mx.Unlock()
				if err != nil {
					_ = utils.WriteDebugLog(fmt.Sprintf("[UUID: %s] w.Broadcast 发送消息失败 失败信息 [ERROR] : %s", uuid, err))
					_ = utils.WriteInfoLog(fmt.Sprintf("[UUID: %s] w.Broadcast 发送消息失败 用户信息 [INFO] : %v", uuid, *client))
					client.IsCloseConn = true
				}
			}

		}
	}
}

// WsClientRegister
// 用户登录
func (w *WsClientManager) WsClientRegister() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("w *WsClientManager WsClientRegister() [ERROR] : %s", err))
			return
		}
		_ = utils.WriteDebugLog(fmt.Sprintf("w *WsClientManager WsClientRegister() 用户登录协程退出"))
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("w *WsClientManager WsClientRegister() 用户登录协程启动"))

	for {
		select {
		case client, ok := <-w.Register:

			if !ok {
				_ = utils.WriteDebugLog(fmt.Sprintf("w.Register close [ERROR] : %s", ok))
				return
			}
			_ = utils.WriteRootLog(fmt.Sprintf("[UUID : %s] : 用户进来了", client.UUID))
			// 判断用户是否存在
			if _, ok := w.Clients[client.UUID]; ok {
				cli := w.Clients[client.UUID]
				// 更新socket连接信息
				cli.Socket = client.Socket
				cli.IpAddress = client.IpAddress
				cli.EndTime = client.EndTime
				cli.IsCloseConn = client.IsCloseConn
			} else {
				w.Clients[client.UUID] = client
			}
		}
	}
}

// WsClientUnRegister
// 用户退出
func (w *WsClientManager) WsClientUnRegister() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("w *WsClientManager WsClientUnRegister() [ERROR] : %s", err))
			return
		}
		_ = utils.WriteDebugLog(fmt.Sprintf("w *WsClientManager WsClientUnRegister() 用户退出协程退出"))
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("w *WsClientManager WsClientUnRegister() 用户退出协程启动"))

	for {
		select {
		case client, ok := <-w.UnRegister:

			if !ok {
				_ = utils.WriteDebugLog(fmt.Sprintf("w.UnRegister close [ERROR] : %s", ok))
				return
			}
			_ = utils.WriteRootLog(fmt.Sprintf("[UUID : %s] : 用户退出了", client.UUID))

			client.IsCloseConn = true
			client.EndTime = time.Now().Unix()
		}
	}
}
