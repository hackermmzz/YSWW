package main

import(
	"time"
	"fmt"
	"strconv"
	"os"
)
//初始化Debug
func DebugInit(){
	//如果需要将debug信息写入文件才进行初始化
	if fileDebug==false{
		return
	}
	//
	debugFile,err:=os.Create("DebugLog.txt")
	if err!=nil{
		fmt.Println("debug初始化失败!")
		return
	}
	//重定向输出流
	os.Stdout=debugFile
	os.Stderr=debugFile
}
// 调试函数
func Debug(message string) {
	now:=time.Now()
	//获取当前时间
	year := strconv.Itoa(now.Year())     
	month := strconv.Itoa(int(now.Month()))  
	day := strconv.Itoa(now.Day())  
	hour:=strconv.Itoa(now.Hour())
	minute:=strconv.Itoa(now.Minute())
	second:=strconv.Itoa(now.Second())
	time_:=year+"-"+month+"-"+day+"-"+hour+"点"+minute+"分"+second+"秒"
	fmt.Println(time_+" "+message)
}