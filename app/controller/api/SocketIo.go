package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"okc/app/service"
	"okc/utils"
	"time"
)

// SocketIO
// socket连接
func SocketIO(c *gin.Context) {
	// 客户uuid
	uuid, _ := c.Cookie("uuid")
	if websocket.IsWebSocketUpgrade(c.Request) {

		upgrader := websocket.Upgrader{CheckOrigin: service.CheckOrigin}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			_ = utils.WriteErrorLog(fmt.Sprintf("upgrader.Upgrade [ERROR] : %s \n", err))
			return
		}
		defer conn.Close()

		// 判断uuid是否存在
		if _, ok := service.WsManager.Clients[uuid]; ok {

			clientInfo := &service.WsClient{
				UUID:        uuid,
				IpAddress:   c.ClientIP(),
				Socket:      conn,
				EndTime:     time.Now().Unix(),
				IsCloseConn: false,
			}

			service.WsManager.Register <- clientInfo

			success, _ := json.Marshal(gin.H{
				"type":    1,
				"message": "success",
				"data": gin.H{
					"uuid": clientInfo.UUID,
				},
			})
			clientInfo.Socket.WriteMessage(websocket.TextMessage, success)

			for {
				_, data, err := clientInfo.Socket.ReadMessage()
				if err != nil {
					log.Printf("clientInfo.Socket.ReadMessage() [ERROR]: %s", err)
					if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						log.Printf("socket close uuid: %s \n", clientInfo.UUID)
						clientInfo.Socket.Close()
						clientInfo.IsCloseConn = true
					}
					return
				}
				log.Printf("有新消息过来了 : %s", string(data))
				log.Printf("service.WsManager.Clients : %v", service.WsManager.Clients)
			}

		} else {

			uuid := utils.GenerateUUID()
			clientInfo := &service.WsClient{
				UUID:        uuid,
				IpAddress:   c.ClientIP(),
				Socket:      conn,
				StartTime:   time.Now().Unix(),
				EndTime:     time.Now().Unix(),
				IsCloseConn: false,
			}

			service.WsManager.Register <- clientInfo

			success, _ := json.Marshal(gin.H{
				"code":    200,
				"message": "success",
				"data": gin.H{
					"uuid": clientInfo.UUID,
				},
			})
			clientInfo.Socket.WriteMessage(websocket.TextMessage, success)

			for {
				_, data, err := clientInfo.Socket.ReadMessage()
				if err != nil {
					log.Printf("clientInfo.Socket.ReadMessage() [ERROR]: %s", err)
					if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						log.Printf("socket close uuid: %s \n", clientInfo.UUID)
						clientInfo.Socket.Close()
						clientInfo.IsCloseConn = true
					}
					return
				}
				log.Printf("有新消息过来了 : %s", string(data))
				log.Printf("service.WsManager.Clients : %v", service.WsManager.Clients)
			}
		}

	}
}

// SocketIOKline
// socket连接
func SocketIOKline(c *gin.Context) {
	// 客户uuid
	uuid, _ := c.Cookie("uuid")
	klineName := c.Query("name")
	if klineName == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage("", utils.ParameterError),
		})
		return
	}

	var isSet bool
	for _, c := range service.CurrencyInstId {
		if c == klineName {
			isSet = true
		}
	}
	if !isSet {
		c.JSON(http.StatusOK, gin.H{
			"type":    404,
			"message": utils.GetLangMessage("", utils.ParameterError),
		})
		return
	}
	if websocket.IsWebSocketUpgrade(c.Request) {

		upgrader := websocket.Upgrader{CheckOrigin: service.CheckOrigin}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			_ = utils.WriteErrorLog(fmt.Sprintf("upgrader.Upgrade [ERROR] : %s \n", err))
			return
		}
		defer conn.Close()

		// 判断uuid是否存在
		if _, ok := service.WsManager.Clients[uuid]; ok {

			clientInfo := &service.WsClient{
				UUID:        uuid,
				IpAddress:   c.ClientIP(),
				Socket:      conn,
				EndTime:     time.Now().Unix(),
				IsCloseConn: false,
			}

			service.WsManager.Clients[clientInfo.UUID] = clientInfo

			success, _ := json.Marshal(gin.H{
				"type":    1,
				"message": "success",
				"data": gin.H{
					"uuid": clientInfo.UUID,
				},
			})
			clientInfo.Socket.WriteMessage(websocket.TextMessage, success)

			go new(service.UserKlineServer).Open(clientInfo, klineName)

			for {
				if clientInfo.IsCloseConn {
					_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s]  客户端退出 \n", clientInfo.UUID))
					return
				}
			}

		} else {

			uuid := utils.GenerateUUID()
			clientInfo := &service.WsClient{
				UUID:        uuid,
				IpAddress:   c.ClientIP(),
				Socket:      conn,
				StartTime:   time.Now().Unix(),
				EndTime:     time.Now().Unix(),
				IsCloseConn: false,
			}

			service.WsManager.Clients[clientInfo.UUID] = clientInfo

			success, _ := json.Marshal(gin.H{
				"code":    200,
				"message": "success",
				"data": gin.H{
					"uuid": clientInfo.UUID,
				},
			})
			clientInfo.Socket.WriteMessage(websocket.TextMessage, success)

			go new(service.UserKlineServer).Open(clientInfo, klineName)

			for {
				if clientInfo.IsCloseConn {
					_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s]  客户端退出 \n", clientInfo.UUID))
					return
				}
			}
		}

	}
}
