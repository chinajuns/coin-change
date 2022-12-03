package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"okc/utils"
	"time"
)

// Tickers
// 行情频道
type Tickers struct {
	instId     string          `json:"inst_id"`  // 产品id
	serverConn *websocket.Conn `json:"conn"`     // socket连接
	IsClose    bool            `json:"is_close"` // 订阅是否关闭
}

// SubOpenTickersServer
// 开启行情订阅(全局)
func (t *Tickers) SubOpenTickersServer(url string) error {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("SubOpenTickersServer() 开启socket订阅 recover() [ERROR] : %s", err))
		}
		t.IsClose = true
		_ = utils.WriteDebugLog(fmt.Sprintf("SubOpenTickersServer() d.IsClose 变更 true 订阅服务出现异常"))
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("SubOpenTickersServer() 开启行情订阅\n"))

	t.IsClose = false
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return err
	}
	t.serverConn = conn

	// 启动心跳机制
	go t.TickersServerPing()
	// 启动检查服务是否关闭
	go t.CheckTickersServerIsClose()

	// 订阅k线
	_ = t.SubOpen()

	for {
		// 检查程序中途是否有关闭现象
		if t.IsClose {
			_ = utils.WriteDebugLog(fmt.Sprintf("SubOpenKlineServer() 服务端异常关闭 \n "))
			return nil
		}

		_, data, err := t.serverConn.ReadMessage()
		if err != nil {
			t.IsClose = true
			_ = utils.WriteErrorLog(fmt.Sprintf("t.conn.ReadMessage() 服务端消息读取失败 [ERROR] : %s \n ", err))
			return err
		}

		//log.Printf("Quotation有新消息进来了: %s", string(data))

		// redis设置延迟数据存储
		//redis := utils.RedisConnect()
		//expireKey := "Quotation-ttl"
		//expireTime := 2
		//ttl, err := redis.Do("ttl", expireKey)
		//if err != nil {
		//	t.IsClose = true
		//	_ = utils.WriteErrorLog(fmt.Sprintf("redis.Do(\"ttl\", \"Kline-ttl\") redis读取失败 [ERROR] : %s \n ", err))
		//	return err
		//}
		//// TODO: 保存数据
		//if ttl.(int64) == -2 {
		//	t.SwitchMessage(data)
		//	redis.Do("SetEx", expireKey, expireTime, "1")
		//}

		//TODO: 保存数据
		t.SwitchMessage(data)

	}
}

// TickersServerPing
// 向服务发送心跳
func (t *Tickers) TickersServerPing() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("t.TickersServerPing() 向服务发送心跳失败 [ERROR] : %s \n ", err))
			t.IsClose = true
		}
		_ = utils.WriteDebugLog(fmt.Sprintf("t.TickersServerPing() d.IsClose 变更 true 订阅服务出现异常"))
	}()
	for {
		time.Sleep(time.Second * 20)
		err := t.serverConn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err != nil {
			_ = utils.WriteErrorLog(fmt.Sprintf("t.conn.ReadMessage() 服务端消息读取失败 [ERROR] : %s \n ", err))
			t.IsClose = true
			return
		}
	}
}

// CheckTickersServerIsClose
// 检查服务是否关闭
func (t *Tickers) CheckTickersServerIsClose() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("t.CheckTickersServerIsClose() 检查服务是否关闭失败 [ERROR] : %s \n ", err))
		}
		_ = utils.WriteDebugLog("t.CheckTickersServerIsClose() 检查服务退出")
	}()

	for {
		time.Sleep(time.Second)
		if t.IsClose {
			// 重新启动一个订阅Kline协程
			go new(Candle1D).SubOpenKlineServer("wss://ws.okx.com:8443/ws/v5/public")
			return
		}
	}
}

// SubOpen
// 开启行情订阅
func (t *Tickers) SubOpen() error {

	for _, instId := range CurrencyInstId {
		_ = utils.WriteDebugLog(fmt.Sprintf("SubOpen() 开启行情订阅 币种: %s \n", instId))
		t.instId = instId
		sendMsg := gin.H{
			"op": "subscribe",
			"args": []interface{}{
				gin.H{
					"channel": "tickers",
					"instId":  t.instId,
				},
			},
		}
		sendMsgJson, _ := json.Marshal(sendMsg)

		t.serverConn.WriteMessage(websocket.TextMessage, sendMsgJson)
	}

	return nil
}

// SwitchMessage
// 判断消息
func (t *Tickers) SwitchMessage(message []byte) {
	var face map[string]interface{}
	_ = json.Unmarshal(message, &face)

	if _, ok := face["event"]; !ok {
		args, ok := face["arg"].(map[string]interface{})
		if !ok {
			//_ = utils.WriteErrorLog(fmt.Sprintf("face[\"arg\"].(map[string]interface{}) [ERROR] : %v", ok))
			return
		}
		data := face["data"].([]interface{})
		mapName := fmt.Sprintf("QHOTATION-%s", args["instId"])
		t.saveCurrencyTickersData(data, mapName)
	}

}

// saveCurrencyTickersData
// 保存币种行情数据
func (t *Tickers) saveCurrencyTickersData(data []interface{}, mapName string) {
	mg := utils.Mongo
	//log.Println("data : ", data)
	_, err := mg.Collection(mapName).InsertMany(context.Background(), data)
	if err != nil {
		_ = utils.WriteErrorLog(fmt.Sprintf("mg.Collection(%s).InsertOne(context.Background(), maps) [ERROR] : %s ", mapName, err))
		return
	}
}

// CurrencyQuotationStruct
type CurrencyQuotationStruct struct {
	Id        string `json:"id" bson:"_id"`              //
	InstId    string `json:"instId" bson:"instId"`       //
	InstType  string `json:"instType" bson:"instType"`   //
	AskPx     string `json:"askPx" bson:"askPx"`         //
	SodUtc8   string `json:"sodUtc8" bson:"sodUtc8"`     //
	Last      string `json:"last" bson:"last"`           //
	AskSz     string `json:"askSz" bson:"askSz"`         //
	BidPx     string `json:"bidPx" bson:"bidPx"`         //
	BidSz     string `json:"bidSz" bson:"bidSz"`         //
	Open24h   string `json:"open24h" bson:"open24h"`     //
	High24h   string `json:"high24h" bson:"high24h"`     //
	VolCcy24h string `json:"volCcy24h" bson:"volCcy24h"` //
	Vol24h    string `json:"vol24h" bson:"vol24h"`       //
	Ts        string `json:"ts" bson:"ts"`               //
	LastSz    string `json:"lastSz" bson:"lastSz"`       //
	Low24h    string `json:"low24h" bson:"low24h"`       //
	SodUtc0   string `json:"sodUtc0" bson:"sodUtc0"`     //
}

// SendCurrencyQuotationStruct
// 发送币种Quotation结构体
type SendCurrencyQuotationStruct struct {
	Code    string                               `json:"code,omitempty"` // 状态码
	Type    string                               `json:"type,omitempty"` // 类型
	Data    map[string][]CurrencyQuotationStruct `json:"data,omitempty"` // 数据
	isClose bool                                 `json:"is_close"`       // 管道是否关闭
}

// SendMessage
// 发送消息
func (s *SendCurrencyQuotationStruct) SendMessage() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyQuotationStruct SendMessage() [ERROR] : %s", err))
			return
		}
		_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyQuotationStruct SendMessage() Quotation广播协程退出"))
		s.isClose = true
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyQuotationStruct SendMessage() Quotation广播协协程启动"))

	go s.CheckSendMessage()

	for {
		time.Sleep(time.Second * 2)
		currency := make(map[string][]CurrencyQuotationStruct)
		mg := utils.Mongo
		// 遍历币种
		for _, instId := range CurrencyInstId {
			container := make([]CurrencyQuotationStruct, 0)
			quotation := new(CurrencyQuotationStruct)

			mapName := fmt.Sprintf("QHOTATION-%s", instId)
			err := mg.Collection(mapName).FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(quotation)
			if err != nil {
				_ = utils.WriteDebugLog(fmt.Sprintf("mg.Collection(%s).Find() [ERROR] : %s", mapName, err))
				return
			}
			container = append(container, *quotation)
			currency[instId] = container
		}

		s.Code = "ok"
		s.Type = "Quotation"
		s.Data = currency
		s.isClose = false
		data, _ := json.Marshal(s)
		WsManager.Broadcast <- data
	}
}

// CheckSendMessage
// 检查发送消息
func (s *SendCurrencyQuotationStruct) CheckSendMessage() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyQuotationStruct CheckSendMessage() [ERROR] : %s", err))
			return
		}
		_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyQuotationStruct CheckSendMessage() 检查Quotation广播协程退出"))
		s.isClose = true
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("s *SendCurrencyQuotationStruct CheckSendMessage() 检查Quotation广播协协程启动"))

	for {
		time.Sleep(time.Second)
		if s.isClose {
			// 启动一个Quotation广播
			go s.SendMessage()
			return
		}
	}
}

// QueryLastQuotationPriceByCurrencyName
// 根据币种名称获取最新行情的价格
func QueryLastQuotationPriceByCurrencyName(currencyName string) (string, error) {
	collation := fmt.Sprintf("QHOTATION-%s", currencyName)
	mongo := utils.Mongo
	quotation := new(CurrencyQuotationStruct)
	err := mongo.Collection(collation).FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(quotation)
	if err != nil {
		return "", err
	}
	return quotation.Last, nil
}
