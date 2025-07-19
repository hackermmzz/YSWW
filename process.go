package main

import (
	"path/filepath"
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"strconv"
)
//
type SocketConn struct {
	name  	string
	conn    net.Conn
	scanner *bufio.Scanner
	lock 	sync.RWMutex
}

var (
	ModelProcess = make(map[string]*SocketConn)
)

func (c *SocketConn) Write(msg ...string) error {
	//提前上锁,防止出现意外
	c.lock.Lock()
	defer c.lock.Unlock()
	//
	for _, arg := range msg {
		_, err := fmt.Fprintf(c.conn, "%s", arg+"\n")
		if err != nil {
			return errors.New("写入失败")
		}
	}
	return nil
}
func (c *SocketConn) Read() (string, error) {
	//读取不需要加互斥锁
	if c.scanner.Scan() {
		return c.scanner.Text(), nil
	}
	return "", c.scanner.Err()
}
func MakeSocketConn(name string, port int) (*SocketConn, error) {
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	//接受客户端连接
	fmt.Println("开始监听端口:"+strconv.Itoa(port))
	conn, err := listener.Accept()
	//
	if err != nil {
		return nil, err
	}
	//
	ret := &SocketConn{
		name:     name,
		conn:    conn,
		scanner: bufio.NewScanner(conn),
	}
	//
	return ret, nil
}
//统一处理管道
func ProcessAIChannel(type_ string,name string,update_fun func(string)error){
	cmd:=ModelProcess[type_]
	//目前只根据generate来定位一条记录
	for{
		generate,err:=cmd.Read()
		if err!=nil{
			Debug(name+"处理管道出现问题!"+err.Error())
		}else{
			err:=update_fun(generate)
			if err!=nil{
				Debug(name+"处理管道在更新数据库记录时出现问题:"+err.Error())
			}else{
				Debug(name+"生成成功:"+generate)
			}
		}
	}
}
//虚拟穿衣处理管道
func ProcessXNCYChannel(){
	ProcessAIChannel("XNCY","虚拟穿衣",UpdateXNCYStatus)
}
//人脸风格化处理管道
func ProcessRLFGHChannel(){
	ProcessAIChannel("RLFGH","人脸风格化",UpdateRLFGHStatus)
}
//文生图
func ProcessWSTChannel(){
	ProcessAIChannel("WST","文生图",UpdateWSTStatus)
}
//人脸肖像化
func ProcessPortraitChannel(){
	ProcessAIChannel("Portrait","人脸肖像化",UpdatePortraitStatus)
}
//智能抠图
func ProcessKouTuChannel(){
	ProcessAIChannel("KouTu","智能抠图",UpdateKouTuStatus)
}
//处理化所有处理管道
func ProcessAIFinishChannel(){
	//虚拟穿衣
	go ProcessXNCYChannel()
	//人脸风格化
	go ProcessRLFGHChannel()
	//文生图
	go ProcessWSTChannel()
	//人脸肖像化
	go ProcessPortraitChannel()
	//智能抠图
	go ProcessKouTuChannel()
}
//初始化所有的AI模型
func ProcessAIModelInit()error{
	g_err:=make([]error,0)
	var wait sync.WaitGroup
	var lock sync.RWMutex
	port:=15263//socket通信端口起始地址
	//配置加载模型数组
	name:=[...]string{"虚拟穿衣","人脸风格化","文生图","人脸肖像化","智能抠图"};
	key:=[...]string{"XNCY","RLFGH","WST","Portrait","KouTu"};
	//
	modelCnt:=len(name)//加载5个模型
	wait.Add(modelCnt)
	//
	for i:=0;i<modelCnt;i+=1{
		go func(i int){
			defer wait.Done()
			cmd,err:=MakeSocketConn(name[i],port+i)
			//
			lock.Lock()
			defer lock.Unlock()
			//
			ModelProcess[key[i]]=cmd
			if err!=nil{
				g_err=append(g_err,err)
			}else{
				Debug(name[i]+"模型加载成功")
			}
		}(i)
	}
	//检查所有的错误
	wait.Wait() 
	if len(g_err)!=0{
		return g_err[0]
	}
	//运行处理管道
	go ProcessAIFinishChannel()
	//
	return nil
}
func ProcessAIModel(model_type string,user string,arguments map[string]interface{}){
	if model_type=="RLFGH"{
		ProcessRlfgh(user,arguments["face"].(string),arguments["wordFile"].(string))
	}else if model_type=="WST"{
		ProcessWST(user,arguments["wordFile"].(string))
	}else if model_type=="XNCY"{
		ProcessXNCY(user,arguments["person"].(string),arguments["clothes"].(string))
	}else if model_type=="Portrait"{
		ProcessPortrait(user,arguments["face"].(string))
	}else if model_type=="KouTu"{
		ProcessKouTu(user,arguments["image"].(string))
	}
}
func ProcessRlfgh(user string ,face string,wordFile string){
	wordFilePath, err := filepath.Abs(SourceDirectory+"/"+wordFile)
	if err != nil {
		Debug("can't get absolute path: " + wordFile)
	}
	faceFilePath,err1 := filepath.Abs(SourceDirectory+"/"+face)
	if err1 != nil {
		Debug("can't get absolute path: " + faceFilePath)
	}
	
	//
	merge:=RandomFileNameWithSuffix(".png")
	//
	cmd :=ModelProcess["RLFGH"] 
	err=cmd.Write(faceFilePath,wordFilePath,SourceDirectory,merge)//写入参数
	if err!=nil{
		Debug("写入人脸风格化管道失败!")
	}
	//向数据库加入一条记录
	err= AddDataToRLFGH(user,face,wordFile,merge)
	if err != nil {
		Debug(err.Error())
	}
}
func ProcessWST(user string, wordFile string) {
	wordFilePath, err := filepath.Abs(SourceDirectory+"/"+wordFile)
	if err != nil {
		Debug("can't get absolute path: " + wordFile)
	}
	
	
	//
	merge:=RandomFileNameWithSuffix(".png")
	//
	cmd :=ModelProcess["WST"] 
	err=cmd.Write(wordFilePath,SourceDirectory,merge)
	if err!=nil{
		Debug("写入文生图管道失败!")
	}
	//
	err = AddDataToWST(user,wordFile,merge)
	if err != nil {
		Debug(err.Error())
	}
}
func ProcessXNCY(user string, person string, clothes string) {
	// 获取绝对路径
	personPath, err := filepath.Abs(SourceDirectory+"/"+person)
	if err != nil {
		Debug("can't get absolute path: " + person)
	}
	clothesPath, err := filepath.Abs(SourceDirectory+"/"+clothes)
	if err != nil {
		Debug("can't get absolute path: " + clothes)
		return
	}
	
	// 执行SourceProcessProgram并获取输出值
	merge := RandomFileNameWithSuffix(".png")
	//
	cmd:=ModelProcess["XNCY"]
	err=cmd.Write(personPath,clothesPath,SourceDirectory,merge)
	if err!=nil{
		Debug("写入虚拟穿衣管道失败!")
	}
	//
	err = AddDataToXNCY(user, person, clothes, merge)
	if err != nil {
		Debug(err.Error())
	}
}
func ProcessPortrait(user string,face string){
	// 获取绝对路径
	facePath, err := filepath.Abs(SourceDirectory+"/"+face)
	if err != nil {
		Debug("can't get absolute path: " + face)
	}
	
	// 执行SourceProcessProgram并获取输出值
	merge := RandomFileNameWithSuffix(".png")

	//
	cmd :=ModelProcess["Portrait"] 
	err=cmd.Write(facePath,SourceDirectory,merge)
	if err!=nil{
		Debug("写入人脸肖像化管道失败!")
	}
	//
	err = AddDataToPortrait(user,face,merge)
	if err != nil {
		Debug(err.Error())
	}
}
func ProcessKouTu(user string,image string){
	// 获取绝对路径
	imagePath, err := filepath.Abs(SourceDirectory+"/"+image)
	if err != nil {
		Debug("can't get absolute path: " + image)
	}
	
	// 执行SourceProcessProgram并获取输出值
	merge := RandomFileNameWithSuffix(".png")

	//
	cmd :=ModelProcess["KouTu"] 
	err=cmd.Write(imagePath,SourceDirectory,merge)
	if err!=nil{
		Debug("写入抠图管道失败!")
	}
	//
	err = AddDataToKouTu(user,image,merge)
	if err != nil {
		Debug(err.Error())
	}
}