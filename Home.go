package main

import(
	"net/http"
	"strconv"
	"errors"
)

//浏览量增加
func ProcessPlayCountHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		Debug("用户不存在!")
		return
	}
	//
	id:=r.Header.Get("id")
	err:=IncreasePlayCountById(id)
	if err!=nil{
		Debug(user +" error increase playcount "+err.Error())
		return
	}
}
//浏览量数量查询
func ProcessPlayCountQueryHanlder(cookie string,w http.ResponseWriter,r*http.Request){
	//
	exist,_:=CheckCookieLegal(cookie)
	if !exist{
		Debug("用户不存在!")
		return
	}
	//
	err:=ProcessPlayCountQuery(cookie,r,w)
	if err!=nil{
		Debug(err.Error())
	}
}
//
func ResponsePlayCountQuery(w http.ResponseWriter,cnt int)error{
	cc:=strconv.Itoa(cnt)
	_, err:= w.Write([]byte(cc))
	return err
}
func ProcessPlayCountQuery(cookie string,r*http.Request,w http.ResponseWriter)error{
	//检查cookie是否合法
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		return errors.New("cookie不存在")
	}
	//
	id:=r.Header.Get("id")
	cnt,err:=QueryPlayCountById(id)
	if err!=nil{
		return errors.New(user +" error query playcount "+err.Error())
	}
	err=ResponsePlayCountQuery(w,cnt)
	if err!=nil{
		return errors.New(user +" error query playcount "+err.Error())
	}
	return nil
}


//
func ResponseHomeInfo(cookie string,w http.ResponseWriter,requestType string) error{
	//
	info:=""
	if requestType=="0"{
		info=HomeInfo
	}else {
		_,exist:=HomeInfoDiv[requestType]
		if exist{
			info=HomeInfoDiv[requestType]
		}
	}
	///////////////////
	if len(info)==0{
		info="{}"
	}
	///////////////////
	_, err:= w.Write([]byte(info))
	return err
}
//主页
func ProcessHomeHandler(cookie string,w http.ResponseWriter,r*http.Request){
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		Debug("用户不存在!")
		return
	}
	//
	type2:= r.Header.Get("type2")
	err := ResponseHomeInfo(cookie, w,type2)
	if err != nil {
		Debug(user + " error uploading Home info " + err.Error())
		return
	}
}

