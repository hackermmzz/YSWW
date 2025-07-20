package main

import(
	"net/http"
	"fmt"
	"errors"
	"unicode"
)
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
	account := r.Header.Get("account")
	passwd:=r.Header.Get("password")
	var json Json;
	cookie,set_cookie:=GenerateCookie(account,passwd)//生成一个cookie用于维护账户的信息
	//检查是否可以登录
	if CheckCanLogin(account,passwd,cookie){
		//将cookie加到维护队列,无论是否之前是否维护了
		AddCookie(account,cookie)
		//
		json.AppendString("msg","登录成功")
		json.AppendBool("status",true)
		//设置一下set-cookie
		w.Header().Set("Set-Cookie",set_cookie)
	}else{
		json.AppendString("msg","账号不存在或者密码错误!")
		json.AppendBool("status",false)
		return errors.New("用户"+account+"不存在")
	}
	_,err:=w.Write([]byte(json.Get()))
	return err
}

//登录
func ProcessLoginHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessLoginRequest(r,w,cookie);
	if err!=nil{
		fmt.Println(err)
	}else{
		Debug("用户:"+r.Header.Get("account")+"登录成功!")
	}
}