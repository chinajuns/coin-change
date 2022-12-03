package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"okc/utils"
	"time"
)

// SubScribeManager
// 订阅管理结构体
type SubScribeManager struct {
	serverConn       *websocket.Conn        // 服务socket连接
	SubScribeChanMap map[string]interface{} // 订阅管道集合
	readMsgChan      chan []byte            // 读消息管道
	writeMsgChan     chan []byte            // 写消息管道
	clientConn       *WsClient
}

// SubScribeOpen
// 开启socket订阅
func (s *SubScribeManager) SubScribeOpen(url string, client *WsClient) error {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("SubScribeOpen() 开启socket订阅 UUID:[%s]  recover() [ERROR] : %s", client.UUID, err))
		}
		client.IsCloseServer = true
	}()
	client.IsCloseServer = false
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return err
	}
	s.serverConn = conn
	s.clientConn = client
	s.SubScribeChanMap = make(map[string]interface{})
	s.readMsgChan = make(chan []byte, 1024)
	s.writeMsgChan = make(chan []byte, 1024)

	defer close(s.readMsgChan)
	defer close(s.writeMsgChan)

	go s.subScribeWriteOpen()
	go s.subScribeReadOpen()

	// 订阅k线
	Candle1D := new(Candle1D)
	_ = Candle1D.SubOpen(s, "BTC-USDT")
	_ = Candle1D.SubOpen(s, "OKB-USDT")
	_ = Candle1D.SubOpen(s, "XRP-USDC")
	_ = Candle1D.SubOpen(s, "SOL-USDT")
	_ = Candle1D.SubOpen(s, "DOGE-USDC")
	_ = Candle1D.SubOpen(s, "OKDOT1-DOT")
	_ = Candle1D.SubOpen(s, "DAI-USDT")
	_ = Candle1D.SubOpen(s, "SHIB-USDT")
	_ = Candle1D.SubOpen(s, "AVAX-ETH")
	_ = Candle1D.SubOpen(s, "WBTC-ETH")
	_ = Candle1D.SubOpen(s, "LEO-USDT")
	_ = Candle1D.SubOpen(s, "LTC-OKB")
	_ = Candle1D.SubOpen(s, "LINK-BTC")

	// TODO: 群发或单独发送
	for {
		_, data, err := s.serverConn.ReadMessage()
		if err != nil {
			_ = utils.WriteErrorLog(fmt.Sprintf("s.conn.ReadMessage() 服务端消息读取失败 [UUID:%s] [ERROR] : %s \n ", s.clientConn.UUID, err))
			return err
		}

		log.Printf("有新消息进来了: %s", string(data))
		// TODO: 保存数据
		Candle1D.SwitchMessage(data)

		// TODO: 单独发送
		err = s.clientConn.Socket.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				s.clientConn.IsCloseConn = true
			}
			_ = utils.WriteErrorLog(fmt.Sprintf("s.clientConn.Socket.WriteMessage() 服务端消息发送失败 [UUID:%s] [ERROR] : %s \n", s.clientConn.UUID, err))
			return err
		}
	}
}

// subScribeReadOpen
// 开启订阅读取消息监听
func (s *SubScribeManager) subScribeReadOpen() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("subScribeReadOpen() 开启订阅读取消息监听 [UUID:%s] recover() [ERROR] : %s \n", s.clientConn.UUID, err))
		}
	}()

	_ = utils.WriteDebugLog(fmt.Sprintf("subScribeReadOpen() [UUID:%s] 开启订阅读取消息监听 \n", s.clientConn.UUID))

	for {
		//time.Sleep(time.Second * 5)
		// TODO: 读取消息写入server服务
		select {
		case msg, ok := <-s.readMsgChan:
			if !ok {
				return
			}

			err := s.serverConn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				_ = utils.WriteErrorLog(fmt.Sprintf("s.serverConn.WriteMessage [UUID:%s] [ERROR] : %s \n", s.clientConn.UUID, err))
			}
		}

	}
}

// subScribeWriteOpen
// 开启订阅写入消息监听
func (s *SubScribeManager) subScribeWriteOpen() {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("subScribeWriteOpen() 开启订阅写入消息监听 [UUID:%s] recover() [ERROR] : %s \n", s.clientConn.UUID, err))
		}
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("subScribeWriteOpen() [UUID:%s] 开启订阅写入消息监听 \n", s.clientConn.UUID))
	for {
		time.Sleep(time.Second * 5)
		// TODO: 读取写入管道消息转发给读消息管道
		select {
		case msg, ok := <-s.writeMsgChan:
			if !ok {
				return
			}

			if len(msg) != 0 {
				log.Println("读取写入管道消息转发给读消息管道")
				s.readMsgChan <- msg
			}

		}
	}

}

// SubScribeWrite
// 订阅消息写入
func (s *SubScribeManager) SubScribeWrite(message []byte) {
	defer func() {
		err := recover()
		if err != nil {
			_ = utils.WriteInfoLog(fmt.Sprintf("SubScribeWrite() 订阅消息写入 [UUID:%s] recover() [ERROR] : %s \n", s.clientConn.UUID, err))
		}
	}()
	_ = utils.WriteDebugLog(fmt.Sprintf("SubScribeWrite() [UUID:%s] 订阅消息写入: %s \n", s.clientConn.UUID, string(message)))

	s.writeMsgChan <- message
}
