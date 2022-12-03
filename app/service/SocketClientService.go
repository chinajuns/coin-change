package service

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// CheckOrigin
// 检查是否跨域
func CheckOrigin(r *http.Request) bool {
	return true
}

// WsClient
// 客户端连接参数
type WsClient struct {
	UUID          string          // 连接uuid
	IpAddress     string          // 连接ip地址
	IpSource      string          // 连接ip来源
	UserId        interface{}     // 用户id
	Socket        *websocket.Conn // socket连接
	Send          chan []byte     // 发送的数据
	StartTime     int64           // 首次连接时间
	EndTime       int64           // 最后一次连接时间
	ExpireTime    int64           // 心跳检测时间
	IsCloseConn   bool            // 连接是否关闭
	IsCloseServer bool            // 服务是否关闭
}
