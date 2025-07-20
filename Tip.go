package main

import(
	"net/http"
	"errors"
)
//判断是否为AI模型生成完成的提醒
func CheckIsAIModelCompleteTip(tip string)bool{
	return tip=="XNCY_SUCCESS"||tip=="WST_SUCCESS"||tip=="RLFGH_SUCCESS"||tip=="Portrait_SUCCESS"||tip=="KouTu_SUCCESS";
}
//处理各种提醒
func ProessTipHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//如果不是内部cookie直接返回
	if !CheckCookieUsedPrivate(cookie){
		return
	}
	//处理各种提醒类型的请求
	tp2:=r.Header.Get("type2")
	if CheckIsAIModelCompleteTip(tp2){
		err:=ProcessAIModelCompleteTip(tp2,w,r)
		if err!=nil{
			Debug(err.Error())
		}
	}else{
		Debug(tp2+"提醒不合法!")
	}
}
//处理AI模型生成完成的提醒
func ProcessAIModelCompleteTip(tp string,w http.ResponseWriter,r*http.Request)error{
	//将响应体解析成json
	msg,err:=ProcessJsonRequestBody(r)
	if err!=nil{
		return err
	}
	//每个请求包必须包含这几个属性
	user:=msg["user"].(string)
	generate:=msg["generate"].(string)
	//
	if tp=="XNCY_SUCCESS"{
		UpdateXNCYStatus(generate)
	}else if tp=="WST_SUCCESS"{
		UpdateWSTStatus(generate)
	}else if tp=="RLFGH_SUCCESS"{
		UpdateRLFGHStatus(generate)
	}else if tp=="Portrait_SUCCESS"{
		UpdatePortraitStatus(generate)
	}else if tp=="KouTu_SUCCESS"{
		UpdateKouTuStatus(generate)
	}else{
		return errors.New("大模型提醒不正确!")
	}
	Debug("用户:"+user+"生成了:"+generate)
	return nil
}