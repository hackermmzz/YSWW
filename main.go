package main

import (
	"math/rand"
	"time"
	"fmt"
)

func main() {
	//随机化
	rand.Seed(time.Now().UnixNano())
	//初始化数据库
	InitDatabase()
	//
	//初始化邮件发送模块
	EmailSendInit()
	//加载大模型
	err:=ProcessAIModelInit()
	if err!=nil{
		fmt.Println(err)
	}
	//配置处理回调函数
	ConfigRequestHandler()
	//启动服务器
	InitServer()
}
