package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"okc/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Candle1D
// K线频道
type Candle1D struct {
	instId     string          `json:"inst_id"`  // 产品id
	serverConn *websocket.Conn `json:"conn"`     // socket连接
	IsClose    bool            `json:"is_close"` // 订阅是否关闭
}

// SubOpenKlineServer
// 开启K线订阅(全局)
func (d *Candle1D) SubOpenKlineServer(url string) error {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("SubOpenKlineServer() 开启socket订阅 recover() [ERROR] : %s", err))
		}
		d.IsClose = true
		_ = utils.WriteDebugLog(fmt.Sprintf("SubOpenKlineServer() d.IsClose 变更 true 订阅服务出现异常"))
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("SubOpenKlineServer() 开启K线订阅\n"))

	d.IsClose = false
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return err
	}
	d.serverConn = conn

	// 启动心跳机制
	go d.KlineServerPing()
	// 启动检查服务是否关闭
	go d.CheckKlineServerIsClose()

	// 订阅k线
	_ = d.SubOpen2()

	for {
		// 检查程序中途是否有关闭现象
		if d.IsClose {
			_ = utils.WriteDebugLog(fmt.Sprintf("SubOpenKlineServer() 服务端异常关闭 \n "))
			return nil
		}

		_, data, err := d.serverConn.ReadMessage()
		if err != nil {
			d.IsClose = true
			_ = utils.WriteErrorLog(fmt.Sprintf("d.conn.ReadMessage() 服务端消息读取失败 [ERROR] : %s \n ", err))
			return err
		}

		//log.Printf("Kline有新消息进来了: %s", string(data))

		// redis设置延迟数据存储
		//redis := utils.RedisConnect()
		//expireKey := "Kline-ttl"
		//expireTime := 2
		//ttl, err := redis.Do("ttl", expireKey)
		//if err != nil {
		//	d.IsClose = true
		//	_ = utils.WriteErrorLog(fmt.Sprintf("redis.Do(\"ttl\", \"Kline-ttl\") redis读取失败 [ERROR] : %s \n ", err))
		//	return err
		//}
		//// TODO: 保存数据
		//if ttl.(int64) == -2 {
		//	d.SwitchMessage(data)
		//	redis.Do("SetEx", expireKey, expireTime, "1")
		//}

		d.SwitchMessage(data)

	}
}

// KlineServerPing
// 向服务发送心跳
func (d *Candle1D) KlineServerPing() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("d.KlineServerPing() 向服务发送心跳失败 [ERROR] : %s \n ", err))
			d.IsClose = true
		}
		_ = utils.WriteDebugLog(fmt.Sprintf("d.KlineServerPing() d.IsClose 变更 true 订阅服务出现异常"))
	}()
	for {
		time.Sleep(time.Second * 20)
		err := d.serverConn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err != nil {
			_ = utils.WriteErrorLog(fmt.Sprintf("d.conn.ReadMessage() 服务端消息读取失败 [ERROR] : %s \n ", err))
			d.IsClose = true
			return
		}
	}
}

// CheckKlineServerIsClose
// 检查服务是否关闭
func (d *Candle1D) CheckKlineServerIsClose() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("d.CheckKlineServerIsClose() 检查服务是否关闭失败 [ERROR] : %s \n ", err))
		}
		_ = utils.WriteDebugLog("d.CheckKlineServerIsClose() 检查服务退出")
	}()

	for {
		time.Sleep(time.Second)
		if d.IsClose {
			// 重新启动一个订阅Kline协程
			go new(Candle1D).SubOpenKlineServer("wss://ws.okx.com:8443/ws/v5/public")
			return
		}
	}
}

// SubOpen2
// 开启K线订阅
func (d *Candle1D) SubOpen2() error {

	for _, instId := range CurrencyInstId {
		_ = utils.WriteDebugLog(fmt.Sprintf("SubOpen()2 开启K线订阅 币种: %s \n", instId))
		d.instId = instId
		sendMsg := gin.H{
			"op": "subscribe",
			"args": []interface{}{
				gin.H{
					"channel": "candle1m",
					"instId":  d.instId,
				},
			},
		}
		sendMsgJson, _ := json.Marshal(sendMsg)

		d.serverConn.WriteMessage(websocket.TextMessage, sendMsgJson)
	}

	return nil
}

// SubOpen
// 开启K线订阅
func (d *Candle1D) SubOpen(s *SubScribeManager, instId string) error {
	if instId == "" {
		return errors.New(utils.GetLangMessage("", utils.ParameterError))
	}

	_ = utils.WriteDebugLog(fmt.Sprintf("SubOpen() [UUID:%s] 开启K线订阅 币种: %s \n", s.clientConn.UUID, instId))

	d.instId = instId
	sendMsg := gin.H{
		"op": "subscribe",
		"args": []interface{}{
			gin.H{
				"channel": "candle1D",
				"instId":  d.instId,
			},
		},
	}
	sendMsgJson, _ := json.Marshal(sendMsg)

	s.SubScribeWrite(sendMsgJson)
	return nil
}

// SwitchMessage
// 判断消息
func (d *Candle1D) SwitchMessage(message []byte) {
	var face map[string]interface{}
	_ = json.Unmarshal(message, &face)

	if _, ok := face["event"]; !ok {
		args, ok := face["arg"].(map[string]interface{})
		if !ok {
			//_ = utils.WriteErrorLog(fmt.Sprintf("face[\"arg\"].(map[string]interface{}) [ERROR] : %v", ok))
			return
		}
		data := face["data"].([]interface{})
		mapName := fmt.Sprintf("KLINE-%s", args["instId"])
		d.saveCurrencyKlineData(data, mapName)
	}

}

// saveCurrencyKlineData
// 保存币种K线数据
func (d *Candle1D) saveCurrencyKlineData(data []interface{}, mapName string) {
	mg := utils.Mongo
	//log.Println("mongo 地址 :", &mg)
	maps := d.interfaceTypeChangToMap(data)
	maps["crt"] = time.Now().UnixMilli()
	_, err := mg.Collection(mapName).InsertOne(context.Background(), maps)
	if err != nil {
		_ = utils.WriteErrorLog(fmt.Sprintf("mg.Collection(%s).InsertOne(context.Background(), maps) [ERROR] : %s ", mapName, err))
		return
	}
}

// interfaceTypeChangToMap
// interface类型转换为map类型
func (d *Candle1D) interfaceTypeChangToMap(param []interface{}) map[string]interface{} {
	maps := make(map[string]interface{})

	for _, v := range param {
		value := v.([]interface{})
		maps["ts"] = value[0]
		maps["o"] = value[1]
		maps["h"] = value[2]
		maps["l"] = value[3]
		maps["c"] = value[4]
		maps["vol"] = value[5]
		maps["volCcy"] = value[6]
	}

	return maps
}

// CurrencyKlineStruct
// 币种K线数据结构体
type CurrencyKlineStruct struct {
	Id     string `json:"_id,omitempty" bson:"_id"`       // objectid
	Ts     string `json:"ts,omitempty" bson:"ts"`         // 开始时间
	O      string `json:"o,omitempty" bson:"o"`           // 开盘价格
	H      string `json:"h,omitempty" bson:"h"`           // 最高价格
	L      string `json:"l,omitempty" bson:"l"`           // 最低价格
	C      string `json:"c,omitempty" bson:"c"`           // 收盘价格
	Vol    string `json:"vol,omitempty" bson:"vol"`       // 交易量，以张为单位,如果是衍生品合约，数值为合约的张数, 如果是币币/币币杠杆，数值为交易货币的数量。
	VolCcy string `json:"volCcy,omitempty" bson:"volCcy"` // 交易量，以张为单位,如果是衍生品合约，数值为合约的张数, 如果是币币/币币杠杆，数值为交易货币的数量。
	Crt    int64  `json:"crt,omitempty" bson:"crt"`       // 创建时间
}

// SendCurrencyKlineStruct
// 发送币种kline结构体
type SendCurrencyKlineStruct struct {
	Code    string                           `json:"code,omitempty"` // 状态码
	Type    string                           `json:"type,omitempty"` // 类型
	Data    map[string][]CurrencyKlineStruct `json:"data,omitempty"` // 数据
	isClose bool                             `json:"is_close"`       // 管道是否关闭
}

// SendMessage
// 发送消息
func (s *SendCurrencyKlineStruct) SendMessage() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyKlineStruct SendMessage() [ERROR] : %s", err))
			return
		}
		_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyKlineStruct SendMessage() Kline广播协程退出"))
		s.isClose = true
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyKlineStruct SendMessage() Kline广播协协程启动"))

	go s.CheckSendMessage()

	for {
		time.Sleep(time.Second * 2)
		currency := make(map[string][]CurrencyKlineStruct)
		mg := utils.Mongo
		// 遍历币种
		for _, instId := range CurrencyInstId {
			container := make([]CurrencyKlineStruct, 0)
			kline := new(CurrencyKlineStruct)

			mapName := fmt.Sprintf("KLINE-%s", instId)
			err := mg.Collection(mapName).FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(kline)
			if err != nil {
				_ = utils.WriteDebugLog(fmt.Sprintf("mg.Collection(%s).Find() [ERROR] : %s", mapName, err))
				return
			}
			container = append(container, *kline)

			// 排序

			currency[instId] = container
		}

		s.Code = "ok"
		s.Type = "Kline"
		s.Data = currency
		s.isClose = false
		data, _ := json.Marshal(s)
		WsManager.Broadcast <- data
	}
}

// CheckSendMessage
// 检查发送消息
func (s *SendCurrencyKlineStruct) CheckSendMessage() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyKlineStruct CheckSendMessage() [ERROR] : %s", err))
			return
		}
		_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyKlineStruct CheckSendMessage() 检查Kline广播协程退出"))
		s.isClose = true
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyKlineStruct CheckSendMessage() 检查Kline广播协协程启动"))

	for {
		time.Sleep(time.Second)
		if s.isClose {
			// 启动一个Kline广播
			go s.SendMessage()
			return
		}
	}
}

// QueryLastKlinePriceByCurrencyName
// 根据币种名称获取最新k线的价格
func QueryLastKlinePriceByCurrencyName(currencyName string) (string, error) {
	collation := fmt.Sprintf("KLINE-%s", currencyName)
	mongo := utils.Mongo
	kline := new(CurrencyKlineStruct)
	err := mongo.Collection(collation).FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(kline)
	if err != nil {
		return "", err
	}
	return kline.C, nil
}

// UserKlineServer
// 用户Kline服务
type UserKlineServer struct {
	client     *WsClient   // socket连接
	readerChan chan []byte // 读管道
	writeChan  chan []byte // 写管道
	klineName  string      // kline名称
}

// Open
// 开启订阅
func (u *UserKlineServer) Open(client *WsClient, klineName string) {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer Open 开启订阅异常 [ERROR] : %s \n", client.UUID, err))
			client.IsCloseConn = true
			return
		}
		_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer Open 开启订阅退出 \n", client.UUID))
		client.IsCloseConn = true
	}()
	_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer 开启订阅 : %s \n", client.UUID, klineName))

	u.client = client
	u.writeChan = make(chan []byte, 1024)
	u.readerChan = make(chan []byte, 1024)
	u.klineName = klineName

	defer close(u.writeChan)
	defer close(u.readerChan)

	go u.ReaderMessage()
	go u.WriteMessage()

	for {
		_, data, err := u.client.Socket.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err) {
				_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] u.Client.Socket.ReadMessage 用户退出 [ERROR] : %s \n", client.UUID, err))
				client.IsCloseConn = true
				return
			}
			_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] u.Client.Socket.ReadMessage 读取消息失败 [ERROR] : %s \n", client.UUID, err))
			client.IsCloseConn = true
			return
		}
		if string(data) == "ping" {
			u.Pong()
		}
	}
}

// ReaderMessage
// 读取消息
func (u *UserKlineServer) ReaderMessage() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer ReaderMessage 读取消息异常 [ERROR] : %s \n", u.client.UUID, err))
			u.client.IsCloseConn = true
			return
		}
		_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer ReaderMessage 读取消息退出 \n", u.client.UUID))
		u.client.IsCloseConn = true
	}()
	_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer WriteMessage 开启读取消息 \n", u.client.UUID))

	for {
		if u.client.IsCloseConn {
			_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer WriteMessage 用户已经退出 \n", u.client.UUID))
			return
		}
		select {
		case msg, ok := <-u.readerChan:
			if !ok {
				//u.IsClose = true
				return
			}
			err := u.client.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				if websocket.IsCloseError(err) {
					_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] u.Client.Socket.WriteMessage 用户退出 [ERROR] : %s  \n", u.client.UUID, err))
					u.client.IsCloseConn = true
					return
				}
				_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] u.Client.Socket.WriteMessage 向用户写入消息失败 [ERROR] : %s  \n", u.client.UUID, err))
				u.client.IsCloseConn = true
			}
			break

		}
	}
}

// WriteMessage
// 写入消息
func (u *UserKlineServer) WriteMessage() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer WriteMessage 写入消息异常 [ERROR] : %s \n", u.client.UUID, err))
			u.client.IsCloseConn = true
			return
		}
		_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer WriteMessage 写入消息退出 \n", u.client.UUID))
		u.client.IsCloseConn = true
	}()
	_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer WriteMessage 开启写入消息 \n", u.client.UUID))

	for {
		time.Sleep(time.Second * 2)
		if u.client.IsCloseConn {
			_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer WriteMessage 用户已经退出 \n", u.client.UUID))
			return
		}

		kline := new(CurrencyKlineStruct)
		mongo := utils.Mongo
		err := mongo.Collection(fmt.Sprintf("KLINE-%s", u.klineName)).FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(kline)
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] mongo.Collection [ERROR] : %s \n", u.client.UUID, err))
			u.client.IsCloseConn = true
			return
		}
		msg := gin.H{
			"code": "ok",
			"type": "Kline",
			"data": gin.H{
				u.klineName: []interface{}{kline},
			},
		}
		data, _ := json.Marshal(msg)
		u.readerChan <- data
	}
}

// Pong
// 心跳回复
func (u *UserKlineServer) Pong() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] *UserKlineServer Pong 心跳回复异常 [ERROR] : %s \n", u.client.UUID, err))
			u.client.IsCloseConn = true
			return
		}
	}()

	if u.client.IsCloseConn {
		return
	}

	err := u.client.Socket.WriteMessage(websocket.TextMessage, []byte("pong"))
	if err != nil {
		if websocket.IsCloseError(err) {
			_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] u.Client.Socket.WriteMessage Pong 用户退出 [ERROR] : %s \n", u.client.UUID, err))
			u.client.IsCloseConn = true
			return
		}
		_ = utils.WriteInfoLog(fmt.Sprintf("[UUID:%s] u.Client.Socket.WriteMessage Pong 向用户写入消息失败 [ERROR] : %s  \n", u.client.UUID, err))
		u.client.IsCloseConn = true
	}
}
