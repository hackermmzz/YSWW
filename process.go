package main

import (
	"net/http"
	"strconv"
)
//
//所有的模型通过socket通信实现生成
var (
	ModelProcess = make(map[string]*SocketConn) //目前不使用
	XNCYQueue		SyncQueue[XNCYInfo]	
	WSTQueue		SyncQueue[WSTInfo]
	RLFGHQueue		SyncQueue[RLFGHInfo]
	PortraitQueue	SyncQueue[PortraitInfo]
	KouTuQueue		SyncQueue[KouTuInfo]
)
/*目前不使用socket面向消息传递,大幅度解耦
//统一处理管道
func ProcessAIChannel(type_ string,name string,update_fun func(string)error){
	channel:=ModelProcess[type_]
	//目前只根据generate来定位一条记录
	for{
		generate,err:=channel.Read()
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
	return
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
	//配置加载模型数组
	name:=[...]string{"虚拟穿衣","人脸风格化","文生图","人脸肖像化","智能抠图"};
	key:=[...]string{"XNCY","RLFGH","WST","Portrait","KouTu"};
	port:=[...]int{XNCYChannelPort,RLFGHChannelPort,WSTChannelPort,PortraitChannelPort,KouTuChanelPort};
	//
	modelCnt:=len(name)//加载5个模型
	wait.Add(modelCnt)
	//
	for i:=0;i<modelCnt;i+=1{
		go func(i int){
			defer wait.Done()
			channel,err:=MakeSocketConn(name[i],port[i])
			//
			lock.Lock()
			defer lock.Unlock()
			//
			ModelProcess[key[i]]=channel
			if err!=nil{
				g_err=append(g_err,err)
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
*/
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

//
func ProcessXNCY(user string, person string, clothes string) {
	var err error
	// 执行SourceProcessProgram并获取输出值
	generate := RandomFileNameWithSuffix(".png")
	//
	/*
	channel:=ModelProcess["XNCY"]
	err=channel.Write(person,clothes,generate)
	if err!=nil{
		Debug("写入虚拟穿衣管道失败!")
	}
	*/
	//写入待处理队列
	var info XNCYInfo
	info.account=user
	info.person=person
	info.clothes=clothes
	info.generate=generate
	XNCYQueue.push(info)
	//
	err = AddDataToXNCY(user, person, clothes, generate)
	if err != nil {
		Debug(err.Error())
	}
}

func ProcessWST(user string, wordFile string) {
	var err error
	//
	generate:=RandomFileNameWithSuffix(".png")
	//
	/*
	channel :=ModelProcess["WST"] 
	err=channel.Write(wordFile,generate)
	if err!=nil{
		Debug("写入文生图管道失败!")
	}
	*/
	//写入待处理队列
	var info WSTInfo
	info.account=user
	info.description=wordFile
	info.generate=generate
	WSTQueue.push(info)
	//
	err = AddDataToWST(user,wordFile,generate)
	if err != nil {
		Debug(err.Error())
	}
}

func ProcessRlfgh(user string ,face string,wordFile string){
	var err error
	//
	generate:=RandomFileNameWithSuffix(".png")
	//
	/*
	channel :=ModelProcess["RLFGH"] 
	err=channel.Write(face,wordFile,generate)//写入参数
	if err!=nil{
		Debug("写入人脸风格化管道失败!")
	}
	*/
	//写入待处理队列
	var info RLFGHInfo
	info.account=user
	info.description=wordFile
	info.face=face
	info.generate=generate
	RLFGHQueue.push(info)
	//向数据库加入一条记录
	err= AddDataToRLFGH(user,face,wordFile,generate)
	if err != nil {
		Debug(err.Error())
	}
}

func ProcessPortrait(user string,face string){
	var err error
	// 执行SourceProcessProgram并获取输出值
	generate := RandomFileNameWithSuffix(".png")

	//
	/*
	channel :=ModelProcess["Portrait"] 
	err=channel.Write(face,generate)
	if err!=nil{
		Debug("写入人脸肖像化管道失败!")
	}
	*/
	//写入待处理队列
	var info PortraitInfo
	info.account=user
	info.person=face
	info.generate=generate
	PortraitQueue.push(info)
	//
	err = AddDataToPortrait(user,face,generate)
	if err != nil {
		Debug(err.Error())
	}
}
func ProcessKouTu(user string,image string){
	var err error
	// 执行SourceProcessProgram并获取输出值
	generate := RandomFileNameWithSuffix(".png")

	//
	/*
	channel :=ModelProcess["KouTu"] 
	err=channel.Write(image,generate)
	if err!=nil{
		Debug("写入抠图管道失败!")
	}
	*/
	//写入待处理队列
	var info KouTuInfo
	info.account=user
	info.image=image
	info.generate=generate
	KouTuQueue.push(info)
	//
	err = AddDataToKouTu(user,image,generate)
	if err != nil {
		Debug(err.Error())
	}
}

//处理模型机对本机的待生成任务的请求
func ProcessAIModelTasksHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//如果不是内部cookie直接返回
	if !CheckCookieUsedPrivate(cookie){
		return
	}
	//获取类型
	type2:=r.Header.Get("type2")
	//
	var js Json
	sendStatus:=false
	if type2=="XNCY"{
		XNCYQueue.hold()
		defer XNCYQueue.release()
		sendStatus=len(XNCYQueue.queue)!=0
		//
		for idx,ele:=range XNCYQueue.queue{
			var js_ Json
			js_.AppendString("user",ele.account)
			js_.AppendString("person",ele.person)
			js_.AppendString("clothes",ele.clothes)
			js_.AppendString("generate",ele.generate)
			js.AppendJson(strconv.Itoa(idx),js_)
		}
		//
		XNCYQueue.clear()
	}else if type2=="WST"{
		WSTQueue.hold()
		defer WSTQueue.release()
		sendStatus=len(WSTQueue.queue)!=0
		//
		for idx,ele:=range WSTQueue.queue{
			var js_ Json
			js_.AppendString("user",ele.account)
			js_.AppendString("description",ele.description)
			js_.AppendString("generate",ele.generate)
			js.AppendJson(strconv.Itoa(idx),js_)
		}
		//
		WSTQueue.clear()
	}else if type2=="RLFGH"{
		RLFGHQueue.hold()
		defer RLFGHQueue.release()
		sendStatus=len(RLFGHQueue.queue)!=0
		//
		for idx,ele:=range RLFGHQueue.queue{
			var js_ Json
			js_.AppendString("user",ele.account)
			js_.AppendString("face",ele.face)
			js_.AppendString("description",ele.description)
			js_.AppendString("generate",ele.generate)
			js.AppendJson(strconv.Itoa(idx),js_)
		}
		//
		RLFGHQueue.clear()
	}else if type2=="Portrait"{
		PortraitQueue.hold()
		defer PortraitQueue.release()
		sendStatus=len(PortraitQueue.queue)!=0
		//
		for idx,ele:=range PortraitQueue.queue{
			var js_ Json
			js_.AppendString("user",ele.account)
			js_.AppendString("person",ele.person)
			js_.AppendString("generate",ele.generate)
			js.AppendJson(strconv.Itoa(idx),js_)
		}
		//
		PortraitQueue.clear()
	}else if type2=="KouTu"{
		KouTuQueue.hold()
		defer KouTuQueue.release()
		sendStatus=len(KouTuQueue.queue)!=0
		//
		for idx,ele:=range KouTuQueue.queue{
			var js_ Json
			js_.AppendString("user",ele.account)
			js_.AppendString("image",ele.image)
			js_.AppendString("generate",ele.generate)
			js.AppendJson(strconv.Itoa(idx),js_)
		}
		//
		KouTuQueue.clear()
	}else{
		Debug("没有对应的任务队列!")
	}
	//
	_,err:=w.Write([]byte(js.Get()))
	if err!=nil{
		Debug("发送模型消息队列失败!")
	}else if sendStatus{
		Debug("发送模型消息队列成功!")
	}
}
