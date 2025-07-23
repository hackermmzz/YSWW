package main

import(
	"os"
	"strconv"
	"strings"
	"bufio"
	"path/filepath"
)
var(
	EmailChannelPort=119
	XNCYChannelPort=120
	WSTChannelPort=121
	RLFGHChannelPort=122
	PortraitChannelPort=123
	KouTuChanelPort=124
	UniverseCookie="20050119"//用于通行的cookie,用途只能是上传和下载文件
)

func ConfigInit(){
	////////////////////读取配置文件
	file, err := os.OpenFile("config.db", os.O_RDONLY, 0644)
	if err != nil {
		Debug("Could not open config.db")
		return
	}
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) != 0 {
			lines = append(lines, scanner.Text())
		}
	}
	////////////////////配置信息
	ConfigInfo = make(map[string]string)
	for _, line := range lines {
		split := strings.Split(line, "=")
		if len(split) == 2 {
			ConfigInfo[split[0]] = split[1]
		}
	}
	/////////////////配置变量
	SourceProcessProgram = ConfigInfo["SourceProcessProgram"]
	/*SourceProcessProgram, err = filepath.Abs(SourceProcessProgram)
	if err != nil {
		Debug("Could not find SourceProcessProgram configuration!")
		return
	}*/
	//
	SourceDirectory = ConfigInfo["SourceDirectory"]
	SourceDirectory, err = filepath.Abs(SourceDirectory)
	if err != nil {
		Debug("Could not find SourceDirectory configuration!")
		return
	}
	DataBaseMaxIdleConns, err = strconv.Atoi(ConfigInfo["DataBaseMaxIdleConns"])
	if err != nil {
		Debug("Could not find DataBaseMaxIdleConns configuration!")
		return
	}
	DataBaseName = ConfigInfo["DataBaseName"]
	DataBaseFileNameMaxLength, err = strconv.Atoi(ConfigInfo["DataBaseFileNameMaxLength"])
	if err != nil {
		Debug("Could not find DataBaseFileNameMaxLength")
		return
	}
	FileMaxSize, err = strconv.Atoi(ConfigInfo["FileMaxSize"])
	if err != nil {
		Debug("Could not find FileMaxSize")
		return
	}
	FileMaxSize <<= 20 //FileMaxSize mb

	RequestMaxProcess,err=strconv.Atoi(ConfigInfo["RequestMaxProcess"])
	if err!=nil{
		Debug("读取最大处理数量失败!")
		return
	}

	fileDebugStr:=ConfigInfo["fileDebug"]
	if fileDebugStr=="true"{
		fileDebug=true
	}else if fileDebugStr=="false"{
		fileDebug=false
	}else{
		fileDebug=false
	}
}