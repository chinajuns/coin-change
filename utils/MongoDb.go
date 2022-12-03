package utils

import (
	"context"
	"fmt"
	"okc/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Mongo *mongo.Database
)

func init() {
	// 设置客户端连接池配置
	m := config.Config().MONGODB.(map[interface{}]interface{})
	poolSize := uint64(500)
	config := options.ClientOptions{
		MaxPoolSize: &poolSize,
		MinPoolSize: &poolSize,
	}
	//co := config.ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", m["DATABASE_USERNAME"], m["DATABASE_PASSWORD"], m["DATABASE_HOST"], m["DATABASE_PORT"]))
	co := config.ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.Background(), co)
	if err != nil {
		panic(fmt.Sprintf("mongo.Connect [ERROR] : %s", err))
		return
	}

	database := client.Database(m["DATABASE_NAME"].(string)).Client()
	err = database.Ping(context.Background(), nil)
	if err != nil {
		panic(fmt.Sprintf("mongo.Ping [ERROR] : %s", err))
		return
	}
	Mongo = database.Database(m["DATABASE_NAME"].(string))

}
