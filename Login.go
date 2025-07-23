package main

import(
	"net/http"
	"errors"
	"unicode"
)
//获取当前用户的密码
func GetUserPassword(account string)(string,error){
	passwd_,exist:=AllUsers.Load(account)
	if !exist{
		return "",errors.New("用户不存在!")
	}
	passwd,_:=passwd_.(string)
	return passwd,nil
}
//检查密码是否相同
func CheckPassWordIsSame(current string,expected string)bool{
	return current==expected;
}

//检查密码是否合法
func CheckPassWordLegal(passwd string)error{
	//必须只能含有数字和字母
	for _,ch:=range passwd{
		if !unicode.IsLetter(ch)&&!unicode.IsDigit(ch){
			return errors.New("密码格式不正确,只能由数字和字母组成!")
		}
	}
	return nil
}
//
func CheckCanLogin(account string,passwd string,cookie string)bool{
	/*目前不开启这个功能
	//如果cookie存在,不管密码账号对错直接可以登录(cookie由后端统一生成)
	if exist,_:=CheckCookieLegal(cookie);exist{
		return true
	}
	*/
	//如果当前维护的账户里面不存在或者密码错误，则不可以登录
	if CheckIsUserExist(account){
		passwd_,_:=AllUsers.Load(account)
		return passwd_==passwd
	}
	return false
}
//判断cookie是否为内部使用
func CheckCookieUsedPrivate(cookie string)bool{
	return cookie==UniverseCookie
}
//检测cookie是否合法,分为两种cookie,一种为内部使用(内部使用的cookie功能受限)，一种为用户使用
func CheckCookieLegal(cookie string)(bool,string){
	//如果cookie为通行cookie
	if CheckCookieUsedPrivate(cookie){
		return true,"UniverseCookie"
	}
	//
	user,exist:=Cookie.Load(cookie)
	if !exist{
		return false,""
	}
	return true,user.(string)
}
//根据账号生成对应的cookie的value
func GenerateCookieValue(user string)string{
	return user;
}
//生成cookie
func GenerateCookie(user string,passwd string)(string,string){
	value:=GenerateCookieValue(user)
	name:="token"
	cookie := &http.Cookie{
        Name:     name,
        Value:    value,
        MaxAge:	  3600*24,//一天
        Path:     "/",
        Secure:   IsHTTPS,           // 仅 HTTPS
        HttpOnly: true,           // 禁止 JavaScript 访问
        SameSite: http.SameSiteNoneMode, // 跨站限制
    }
    return 	name+"="+value,cookie.String()
}
//添加到全局cookie
func AddCookie(account string,cookie string)bool{
	_,exist:=Cookie.Load(cookie)
	if !exist{
		Cookie.Store(cookie,account)
		return true
	}
	return false
}
func CheckIsUserExist(account string)bool{
	_,exist:=AllUsers.Load(account)
	return exist
}

//


func ProcessLoginRequest(r*http.Request,w http.ResponseWriter,cookie string)error{
	//解析请求体获取数据
	data,err:=ProcessJsonRequestBody(r)
	if err!=nil{
		return err
	}
	//
	var account,passwd string
	{
		var legal0,legal1 bool
		account,legal0=data["account"].(string)
		passwd,legal1=data["password"].(string)
		if passwd==""||account==""||!legal0||!legal1{
			return errors.New("存在字段为空!")
		}
	}
	//
	var json Json;
	cookie,set_cookie:=GenerateCookie(account,passwd)//生成一个cookie用于维护账户的信息
	//检查是否可以登录
	if CheckCanLogin(account,passwd,cookie){
		//将cookie加到维护队列,无论是否之前是否维护了
		AddCookie(account,cookie)
		//获取用户昵称
		userinfo,err:=GetUserInfo(account)
		if err!=nil{
			return errors.New("获取用户昵称失败"+err.Error())
		}
		//
		json.AppendBool("vip",bool(userinfo.vip==1))
		json.AppendString("date",userinfo.date)
		json.AppendString("userID",userinfo.userID)
		json.AppendString("msg","登录成功")
		json.AppendBool("status",true)
		json.AppendString("avatar",userinfo.avatar)
		//设置一下set-cookie
		w.Header().Set("Set-Cookie",set_cookie)
	}else{
		json.AppendString("msg","账号不存在或者密码错误!")
		json.AppendBool("status",false)
	}
	_,err=w.Write([]byte(json.Get()))
	return err
}

//登录
func ProcessLoginHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessLoginRequest(r,w,cookie);
	if err!=nil{
		Debug(err.Error())
	}else{
		Debug("用户:"+r.Header.Get("account")+"登录成功!")
	}
}
//删除账户
func ProcessDeleteAccountHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		Debug("cookie不存在")
		return
	}
	//
	data,err:=ProcessJsonRequestBody(r)
	if err!=nil{
		Debug(err.Error())
		return
	}
	//获取密码来确认
	var js Json
	passwd,err:=GetUserPassword(user)
	comfired_passwd,legal:=data["password"].(string)
	if err!=nil||!legal{
		Debug(user+"删除账号异常!")
		js.AppendBool("status",false)
		js.AppendString("msg","服务器异常!")
	}else if !CheckPassWordIsSame(passwd,comfired_passwd){//检查密码是否匹配
		Debug(user+"删除账号错误!")
		js.AppendBool("status",false)
		js.AppendString("msg","密码不匹配!")
	}else if err=DeleteAccount(user);err!=nil{//将账户从数据库删除
		Debug(user+"删除账号异常!")
		js.AppendBool("status",false)
		js.AppendString("msg","服务器异常!")
	}else{
		js.AppendBool("status",true)
		js.AppendString("msg","注销成功!期待下次再见"+user)
		//将信息从维护的全局账户移除
		DeleteUserAllInfoMantancedByProcess(user)
		//
		Debug("用户:"+user+"注销成功")
	}
	//
	_,err=w.Write([]byte(js.Get()))
	
}
//移除所有由进程维护的信息
func DeleteUserAllInfoMantancedByProcess(account string){
	//移除用户信息
	AllUsers.Delete(account)
	//移除cookie信息
	//移除各种ai所需要的信息
	XNCYDataMap.Delete(account)
	RlfghDataMap.Delete(account)
}