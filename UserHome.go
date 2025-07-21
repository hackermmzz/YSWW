package main

import(
	"net/http"
)

//处理用户主页请求
func ProcessUserHomeHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//
	legal,user:=CheckCookieLegal(cookie)
	if !legal{
		Debug("cookie不合法")
		return
	}
	//获取用户信息
	info,err:=GetUserInfo(user)
	if err!=nil{
		Debug("获取用户:"+user+"信息失败!"+err.Error())
		RespondNCK(w)
		return
	}
	var vipinfo VIPInfo
	//如果用户会员已经开通,那么获取会员信息
	if info.vip!=0{
		vipinfo,err=GetVIPInfo(user)
		if err!=nil{
			Debug("获取会员信息出错!"+err.Error())
			RespondNCK(w)
			return
		}
	}
	//打包响应
	var js Json
	js.AppendBool("status",true)
	js.AppendString("userID",info.userID)
	js.AppendString("avatar",info.avatar)
	js.AppendString("regist_date",info.date)
	js.AppendBool("vip",bool(info.vip!=0))
	js.AppendString("email",info.email)
	if info.vip!=0{
		js.AppendString("vip_begin_date",vipinfo.start_time)
		js.AppendString("vip_end_date",vipinfo.end_time)
	}
	//
	_,err=w.Write([]byte(js.Get()))
	if err!=nil{
		Debug(err.Error())
	}
}


//
func ProcessUserIDChangeHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//
	legal,user:=CheckCookieLegal(cookie)
	if !legal{
		Debug("cookie不合法")
		return
	}
	//
	js,err:=ProcessJsonRequestBody(r)
	if err!=nil{
		Debug(err.Error())
		RespondNCK(w)
		return
	}
	//
	userID,legal:=js["userID"].(string)
	if !legal||!CheckUserIDLegal(userID){
		Debug("昵称不合法!")
		RespondNCK(w)
		return
	}
	err=UpdateUserID(user,userID)
	if err!=nil{
		RespondNCK(w)
		Debug(err.Error())
		return
	}
	//
	RespondACK(w)
}

//修改头像
func ProcessAvatarChangeHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//
	legal,user:=CheckCookieLegal(cookie)
	if !legal{
		Debug("cookie不合法")
		return
	}
	//
	fileName:=DownLoadFile(r)
	if fileName==""{
		RespondNCK(w)
		return	
	}
	//更新数据库
	err:=UpdateUserAvatar(user,fileName)
	if err!=nil{
		Debug("更新头像失败")
		RespondNCK(w)
		return
	}
	//
	var js Json
	js.AppendBool("status",true)
	js.AppendString("new_avatar",fileName)
	_,err=w.Write([]byte(js.Get()))
	if err!=nil{
		Debug(err.Error())
	}else{
		Debug("用户:"+user+"修改头像成功!")
	}
}