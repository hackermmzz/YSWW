package main

import(
	"net/http"
)

//修改密码
func ProcessPassWordChangeHandler(cookie string,w http.ResponseWriter,r*http.Request){
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		Debug("cookie不存在")
		return
	}
	//
	var js Json
	//
	new_passwd:=r.Header.Get("new_passwd")//新密码
	origin_Passwd:=r.Header.Get("password")//用户传来的原密码
	current_passwd,_:=GetUserPassword(user)//真正的密码
	if origin_Passwd!=current_passwd{
		js.AppendBool("status",false)
		js.AppendString("msg","原密码错误!")
	}else if CheckPassWordIsSame(current_passwd,new_passwd){
		//如果与之前的密码相同,或者检测密码是否不合法，那么不理睬
		js.AppendBool("status",false)
		js.AppendString("msg","新密码与原来的密码相同!")
	}else if err:=CheckPassWordLegal(new_passwd);err!=nil{
		js.AppendBool("status",true)
		js.AppendString("msg",err.Error())
	}else{
		//修改密码
		err:=ChangePassWord(user,new_passwd)
		//修改对应记录里面对的密码
		AllUsers.Store(user,new_passwd)
		//
		if err!=nil{
			Debug(err.Error())
		}else{
			Debug("用户:"+user+"修改密码成功!")
			js.AppendBool("status",true)
			js.AppendString("msg","修改成功!")
		}
	}
	_,err:=w.Write([]byte(js.Get()))
	if err!=nil{
		Debug(err.Error())
	}
}
