package main

import(
	"net/http"
	"errors"
	"os"
	"encoding/json"
	"io/ioutil"
	"sync"
)
////////////////////////////////////////////////
var XNCYDataMap 	sync.Map
var RlfghDataMap	sync.Map
////////////////////////////////////////////////
func ProcessRlfghRequest(cookie string,r* http.Request,w http.ResponseWriter)error{
	//检查cookie是否合法
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		return errors.New("cookie不存在")
	}
	//
	info_,exist:=RlfghDataMap.Load(user)
	if !exist{
		RlfghDataMap.Store(user,make([]string,2))
		info_,exist=RlfghDataMap.Load(user)
	}
	info:=info_.([]string)
	//
	type2:= r.Header.Get("type2")
	if type2=="Face"{
		face:=DownLoadFile(r)
		if face==""{
			return errors.New("下载图片文件失败")
		}
		info[0]=face
	}else if type2=="Description"{
		wordFile:=ProcessRnfghWordFile(user,r,w)
		if wordFile==""{
			return	errors.New("下载描述文件失败")
		}
		info[1]=wordFile
	}
	if (len(info[1])!=0) && (len(info[0])!=0){
		arguments:=make(map[string]interface{})
		arguments["face"]=info[0]
		arguments["wordFile"]=info[1]
		go ProcessAIModel("RLFGH",user,arguments)
		info[1]=""
		info[0]=""
		Debug("人脸风格化正在生成...")
	}
	return nil
}
func ProcessXNCYRequest(r *http.Request,cookie string,requestType string)error{
		//检查cookie是否合法
		exist,user:=CheckCookieLegal(cookie)
		if !exist{
			return errors.New("cookie不存在")
		}
		fileName := DownLoadFile(r)
		if fileName == "" {
			return errors.New("download file error!")
		}
		
		//
		info_,exist:=XNCYDataMap.Load(user)
		if !exist{
			XNCYDataMap.Store(user,make([]string,2))
			info_,exist=XNCYDataMap.Load(user)
		}
		info:=info_.([]string)
		//
		if(requestType=="UploadPerson"){
			info[0]=fileName
		}else if(requestType=="UploadClothes"){
			info[1]=fileName
		}
		if (len(info[1])!=0)&&(len(info[0])!=0){
			arguments:=make(map[string]interface{})
			arguments["person"]=info[0]
			arguments["clothes"]=info[1]
			go ProcessAIModel("XNCY",user,arguments)
			info[1]=""
			info[0]=""
			Debug("虚拟穿衣正在生成...")
		}
		return nil
}

func ProcessPortraitRequest(cookie string,r*http.Request,w http.ResponseWriter)error{
		//检查cookie是否合法
		exist,user:=CheckCookieLegal(cookie)
		if !exist{
			return errors.New("cookie不存在")
		}
		//
		face:=DownLoadFile(r)
		if face==""{
			return errors.New("文件下载失败")
		}
		arguments:=make(map[string]interface{})
		arguments["face"]=face
		go ProcessAIModel("Portrait",user,arguments)
		return nil;
}
func ProcessWSTRequest(cookie string,r*http.Request,w http.ResponseWriter)error{
		//检查cookie是否合法
		exist,user:=CheckCookieLegal(cookie)
		if !exist{
			return errors.New("cookie不存在")
		}
		//
		wordFileName := ProcessWenShengTu(user, r, w)
		if wordFileName == "" {
			return errors.New("文件下载失败")
		}
		arguments:=make(map[string]interface{})
		arguments["wordFile"]=wordFileName
		go ProcessAIModel("WST",user, arguments)
		return nil
}
func ProcessKouTuRequest(cookie string,r*http.Request,w http.ResponseWriter)error{
		//检查cookie是否合法
		exist,user:=CheckCookieLegal(cookie)
		if !exist{
			return errors.New("cookie不存在")
		}
		//
		image:=DownLoadFile(r)
		if image==""{
			return errors.New("文件下载失败!")
		}
		arguments:=make(map[string]interface{})
		arguments["image"]=image
		go ProcessAIModel("KouTu",user,arguments)
		return nil
}




func ProcessRnfghWordFile(user string,r *http.Request, w http.ResponseWriter) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	wordFileName := RandomFileName()
	localFileName := SourceDirectory + "/" + wordFileName
	out, err := os.Create(localFileName)
	if err != nil {
		Debug("failed to open the file "+localFileName+" for writing")
		return ""
	}
	defer out.Close()
	res:=make(map[string]interface{})
	err = json.Unmarshal(body, &res)
	if err != nil {
		Debug(err.Error())
		return ""
	}
	out.Write([]byte(res["description"].(string)))
	return wordFileName
}
func ProcessWenShengTu(user string, r *http.Request, w http.ResponseWriter) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	wordFileName := RandomFileName()
	localFileName := SourceDirectory + "/" + wordFileName
	out, err := os.Create(localFileName)
	if err != nil {
		Debug("failed to open the file "+localFileName+" for writing\n")
		return ""
	}
	defer out.Close()
	res:=make(map[string]interface{})
	err = json.Unmarshal(body, &res)
	if err != nil {
		Debug(err.Error())
		return ""
	}
	out.Write([]byte(res["description"].(string)))
	return wordFileName
}

//抠图
func ProcessKouTuHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessKouTuRequest(cookie,r,w)
	if err!=nil{
		Debug(err.Error())
	}else{
		Debug("智能抠图正在生成...")
	}
}

//人脸肖像
func ProcessPortraitHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessPortraitRequest(cookie,r,w)
	if err!=nil{
		Debug(err.Error())
	}else{
		Debug("人脸肖像正在生成")
	}
}

//文生图
func ProcessWenShengTuHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessWSTRequest(cookie,r,w)
	if err!=nil{
		Debug(err.Error())
	}else{
		Debug("文生图正在生成...")
	}
}
//人脸风格化
func ProcessRlfghHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessRlfghRequest(cookie,r,w)
	if err!=nil{
		Debug(err.Error())
	}
}

//虚拟穿衣处理操作
func ProcessXNCYHandler(cookie string,w http.ResponseWriter,r*http.Request){
	requestType:=r.Header.Get("type")
	err:=ProcessXNCYRequest(r,cookie,requestType)
	if err!=nil{
		Debug("虚拟穿衣出现异常! "+err.Error())
	}
}
