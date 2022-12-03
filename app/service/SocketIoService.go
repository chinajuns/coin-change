package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"okc/utils"
	"time"
)

var Server *socketio.Server
var SocketIoManage ServerManage

type ServerEvents struct {
}

type ServerManage struct {
	Clients    map[string]*ServerClient
	Register   chan *ServerClient
	UnRegister chan *ServerClient
}

type ServerClient struct {
	Id             string `json:"id"`               // socket连接id
	RegisterTime   int64  `json:"register_time"`    // 登录时间
	UnRegisterTime int64  `json:"un_register_time"` // 退出时间
	RemoteAddr     string `json:"remote_addr"`      // 登录地址
	IsClose        bool   `json:"is_close"`         // 是否退出
}

func init() {
	SocketIoManage.Clients = make(map[string]*ServerClient)
	SocketIoManage.Register = make(chan *ServerClient, 1024)
	SocketIoManage.UnRegister = make(chan *ServerClient, 1024)

	ServerEvents := new(ServerEvents)
	Server = socketio.NewServer(ServerEvents.GenerateConfig())
	Server.OnConnect("/", ServerEvents.Connect)
	Server.OnEvent("/", "Kline", ServerEvents.Kline)
	Server.OnEvent("/", "Quotation", ServerEvents.Quotation)
	Server.OnDisconnect("/", ServerEvents.Disconnect)
}

func (e *ServerEvents) GenerateConfig() *engineio.Options {
	return &engineio.Options{
		PingTimeout:  7 * time.Second,
		PingInterval: 5 * time.Second,
		Transports: []transport.Transport{
			&polling.Transport{
				Client: &http.Client{
					Timeout: time.Minute,
				},
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	}
}

func (e *ServerEvents) Connect(s socketio.Conn) error {
	s.SetContext("")
	client := &ServerClient{
		Id:           s.ID(),
		RegisterTime: time.Now().Unix(),
		RemoteAddr:   s.RemoteAddr().String(),
		IsClose:      false,
	}
	SocketIoManage.Register <- client
	return nil
}

func (e *ServerEvents) Disconnect(s socketio.Conn, reason string) {
	client := SocketIoManage.Clients[s.ID()]
	SocketIoManage.UnRegister <- client
}

func (e *ServerEvents) Kline(s socketio.Conn, msg string) {
	url := s.URL()
	uuid := url.Query().Get("uuid")
	name := url.Query().Get("name")
	lang := url.Query().Get("lang")
	eventName := "Kline"
	if name == "" {
		errMsg := gin.H{
			"code":    500,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		}
		s.Emit(eventName, errMsg)
		return
	}
	var isSet bool
	for _, c := range CurrencyInstId {
		if c == name {
			isSet = true
		}
	}
	if !isSet {
		errMsg := gin.H{
			"code":    500,
			"message": utils.GetLangMessage(lang, utils.CurrencyFindError),
		}
		s.Emit(eventName, errMsg)
		return
	}
	if uuid == "" {
		uuid = utils.GenerateUUID()
	}
	connectMsg := gin.H{
		"code": "ok",
		"type": "login",
		"data": gin.H{
			"uuid": uuid,
		},
		"message": utils.GetLangMessage(lang, utils.Success),
	}
	s.Emit(eventName, connectMsg)

	go func(s socketio.Conn) {
		defer func() {
			err := recover()
			if err != nil {
				_ = utils.WriteInfoLog(fmt.Sprintf("订阅Kline异常协程退出 连接id: %s [ERROR] : %s", s.ID(), err))
				return
			}
			_ = utils.WriteInfoLog(fmt.Sprintf("订阅Kline正常协程退出 连接id: %s", s.ID()))
			return
		}()

		if _, ok := SocketIoManage.Clients[s.ID()]; !ok {
			return
		}

		for {
			client := SocketIoManage.Clients[s.ID()]
			if client.IsClose {
				delete(SocketIoManage.Clients, s.ID())
				return
			}
			time.Sleep(time.Second / 2)

			kline := new(CurrencyKlineStruct)
			mongo := utils.Mongo
			collectionName := fmt.Sprintf("KLINE-%s", name)
			err := mongo.Collection(collectionName).FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(kline)
			if err != nil {
				errMsg := gin.H{
					"code":    500,
					"message": err.Error(),
				}
				s.Emit(eventName, errMsg)
				return
			}
			data := gin.H{
				"code": "ok",
				"type": "Kline",
				"data": kline,
			}
			s.Emit(eventName, data)
		}
	}(s)
}

func (e *ServerEvents) Quotation(s socketio.Conn, msg string) {
	url := s.URL()
	uuid := url.Query().Get("uuid")
	lang := url.Query().Get("lang")
	eventName := "Quotation"

	if uuid == "" {
		uuid = utils.GenerateUUID()
	}
	connectMsg := gin.H{
		"code": "ok",
		"type": "login",
		"data": gin.H{
			"uuid": uuid,
		},
		"message": utils.GetLangMessage(lang, utils.Success),
	}
	s.Emit(eventName, connectMsg)

	go func(s socketio.Conn) {
		defer func() {
			err := recover()
			if err != nil {
				_ = utils.WriteInfoLog(fmt.Sprintf("订阅Quotation异常协程退出 连接id: %s [ERROR] : %s", s.ID(), err))
				return
			}
			_ = utils.WriteInfoLog(fmt.Sprintf("订阅Quotation正常协程退出 连接id: %s", s.ID()))
			return
		}()

		if _, ok := SocketIoManage.Clients[s.ID()]; !ok {
			return
		}

		for {
			client := SocketIoManage.Clients[s.ID()]
			if client.IsClose {
				delete(SocketIoManage.Clients, s.ID())
				return
			}
			time.Sleep(time.Second / 2)
			container := make([]CurrencyQuotationStruct, 0)
			mongo := utils.Mongo
			for _, instId := range CurrencyInstId {
				quotation := new(CurrencyQuotationStruct)
				collectionName := fmt.Sprintf("QHOTATION-%s", instId)
				err := mongo.Collection(collectionName).FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(quotation)
				if err != nil {
					errMsg := gin.H{
						"code":    500,
						"message": err.Error(),
					}
					s.Emit(eventName, errMsg)
					return
				}
				container = append(container, *quotation)
			}

			data := gin.H{
				"code": "ok",
				"type": "Quotation",
				"data": container,
			}

			s.Emit(eventName, data)
		}
	}(s)
}

func (m ServerManage) RegisterGo() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("RegisterGo 用户登录协程异常退出 [ERROR] : %s ", err))
			return
		}
		_ = utils.WriteInfoLog(fmt.Sprintf("RegisterGo 用户登录协程正常退出"))
		return
	}()
	_ = utils.WriteInfoLog(fmt.Sprintf("RegisterGo 开启用户登录协程"))

	for {
		select {
		case client, ok := <-m.Register:
			if !ok {
				_ = utils.WriteInfoLog(fmt.Sprintf("RegisterGo 用户登录协程管道关闭"))
				return
			}
			m.Clients[client.Id] = client
			break
		}
	}
}

func (m ServerManage) UnRegisterGo() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("UnRegisterGo 用户退出协程异常退出 [ERROR] : %s", err))
			return
		}
		_ = utils.WriteInfoLog(fmt.Sprintf("UnRegisterGo 用户退出协程正常退出"))
		return
	}()
	_ = utils.WriteInfoLog(fmt.Sprintf("UnRegisterGo 开启用户退出协程"))

	for {
		select {
		case client, ok := <-m.UnRegister:
			if !ok {
				_ = utils.WriteInfoLog(fmt.Sprintf("UnRegisterGo 用户退出协程管道关闭"))
				return
			}
			client.IsClose = true
			client.UnRegisterTime = time.Now().Unix()
			break
		}
	}
}
