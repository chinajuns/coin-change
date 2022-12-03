package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"okc/app/service"
	"okc/utils"
	"sort"
	"strconv"
)

// KlineHistory
// k线历史记录
func KlineHistory(c *gin.Context) {
	lang := c.GetHeader("lang")
	times := c.PostForm("t")
	instId := c.PostForm("instId")
	limit := c.PostForm("limit")
	bar := c.PostForm("bar")

	if times == "" || instId == "" || instId == "" || bar == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	timesParseInt, _ := strconv.ParseInt(times, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)
	barUint := utils.TimeStrChangeSecond(bar)

	mongo := utils.Mongo
	collectionName := fmt.Sprintf("KLINE-%s", instId)
	collection := mongo.Collection(collectionName)
	log.Println("timesParseInt : ", timesParseInt)
	cur, err := collection.Find(context.Background(), bson.M{
		"crt": bson.M{"$lte": timesParseInt},
	}, options.Find().SetSort(bson.M{"_id": -1}))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"type":    500,
			"message": utils.GetLangMessage(lang, utils.ServerError),
		})
		_ = utils.WriteErrorLog(fmt.Sprintf("collection.Find [ERROR] : %s", err))
		return
	}
	defer cur.Close(context.Background())
	container := make([]service.CurrencyKlineStruct, 0)
	sum := 0
	for cur.Next(context.Background()) {
		var lock = false

		if sum >= int(limitParseInt) {
			break
		}

		kline := new(service.CurrencyKlineStruct)
		cur.Decode(kline)

		if kline.Crt <= timesParseInt {

			if lock == false {
				timesParseInt = kline.Crt
				lock = true
			}

			container = append(container, *kline)
			timesParseInt -= int64(barUint) * 1000
			sum++
		}

	}
	sort.Slice(container, func(i, j int) bool {
		if container[i].Crt < container[j].Crt {
			return true
		}
		return false
	})
	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    container,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}

// KlineHistoryV2
// k线历史记录V2
func KlineHistoryV2(c *gin.Context) {
	lang := c.GetHeader("lang")
	times := c.PostForm("t")
	instId := c.PostForm("instId")
	limit := c.PostForm("limit")
	bar := c.PostForm("bar")

	if times == "" || instId == "" || instId == "" || bar == "" {
		c.JSON(http.StatusOK, gin.H{
			"type":    400,
			"message": utils.GetLangMessage(lang, utils.ParameterError),
		})
		return
	}

	timesParseInt, _ := strconv.ParseInt(times, 10, 64)
	limitParseInt, _ := strconv.ParseInt(limit, 10, 64)
	barUint := utils.TimeStrChangeSecond(bar)

	mongo := utils.Mongo
	collectionName := fmt.Sprintf("KLINE-%s", instId)
	collection := mongo.Collection(collectionName)
	container := make([]service.CurrencyKlineStruct, 0)

	// 按照limit 进行循环
	for i := 0; i < int(limitParseInt); i++ {
		timesParseStr := strconv.Itoa(int(timesParseInt))
		// 游标
		cur, err := collection.Find(context.Background(), bson.M{
			"ts": timesParseStr,
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"type":    500,
				"message": utils.GetLangMessage(lang, utils.ServerError),
			})
			_ = utils.WriteErrorLog(fmt.Sprintf("collection.Find [ERROR] : %s", err))
			return
		}

		klines := make([]*service.CurrencyKlineStruct, 0)

		for cur.Next(context.Background()) {

			kline := new(service.CurrencyKlineStruct)
			cur.Decode(kline)
			klines = append(klines, kline)
		}

		if len(klines) == 0 {
			break
		}

		klineStartData := klines[0]

		klineEndData := klines[len(klines)-1]

		data := service.CurrencyKlineStruct{
			Id:     klineEndData.Id,
			Ts:     klineEndData.Ts,
			O:      klineStartData.O,
			H:      klineEndData.H,
			L:      klineEndData.L,
			C:      klineEndData.C,
			Vol:    klineEndData.Vol,
			VolCcy: klineEndData.VolCcy,
			Crt:    klineEndData.Crt,
		}

		container = append(container, data)

		timesParseInt -= int64(barUint) * 1000
	}

	sort.Slice(container, func(i, j int) bool {
		if container[i].Crt < container[j].Crt {
			return true
		}
		return false
	})
	c.JSON(http.StatusOK, gin.H{
		"type":    "ok",
		"data":    container,
		"message": utils.GetLangMessage(lang, utils.Success),
	})
	return
}
