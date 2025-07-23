package main

import(
	"sync"
	"fmt"
	"net/http"
	"time"
	"strconv"
	"errors"
	"regexp"
)
///////////////////////////////////////
var RegistVerifyCodeMap sync.Map
//////////////////////////////////////
//检测用户是否合法
func CheckAccountLegal(account string)error{
	//长度不能超过20并且格式为A@B.C,只能为字母和数字
	if len(account)>30{
		return errors.New("邮箱长度不能超过30!")
	}else{
		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		match:=re.MatchString(account)
		if !match{
			return errors.New("邮箱格式不正确!")
		}
	}
	return nil
}
//检查userID是否合法
func CheckUserIDLegal(userID string)bool{
	length:=len(userID)
	return length>0&&length<=20;
}
//
func RegistAccount(account string,userID string,password string)error{
	//检测用户和密码是否合法
	err:=CheckAccountLegal(account)
	if err!=nil{
		return err
	}
	err=CheckPassWordLegal(password)
	if err!=nil{
		return err
	}
	//检测用户的昵称是否合法
	if !CheckUserIDLegal(userID){
		return errors.New("用户昵称不合法!")
	}
	//向数据库添加用户
	err=AddUser(account,userID,password,0)
	if err==nil{
		//添加用户到进程
		AllUsers.Store(account,password)
	}
	return err
}

//生成验证码
func GenerateRegistVerifyCode(account string)string{
	t:=time.Now().UnixNano()/int64(1e6)
	str:=strconv.Itoa(int(t))
	length:=len(str)
	return str[length-5:length];
}
//检查注册的用户和验证码是否匹配
func CheckRegistAccountMatchVerifyCode(account string,code string)bool{
	//获取所有的该账号的记录
	record,err:=QueryAccountRegistCode(account)
	if err!=nil{
		Debug(err.Error())
		return false
	}
	for _,info :=range record{
		if info.code==code{
			//将该项设为已过期
			err:=UpdateAccountRegistCodeStatus(account,code)
			if err!=nil{
				return false
			}else{
				return true
			}
		}
	}
	return false;
}
//
func ProcessRegistRequest(r* http.Request,w http.ResponseWriter)error{
	//从请求体拿到数据
	data,err:=ProcessJsonRequestBody(r)
	if err!=nil{
		return err
	}
	//拿到各个字段
	var account,userID,passwd,verifyCode string
	{
		var legal0,legal1,legal2,legal3 bool
		account,legal0= data["account"].(string)
		userID,legal1=data["userID"].(string)
		passwd,legal2=data["password"].(string)
		verifyCode,legal3=data["verifyCode"].(string)
		if account==""||userID==""||passwd==""||verifyCode==""||!legal0||!legal1||!legal2||!legal3{
			return errors.New("存在字段为空!")
		}
	}
	//
	var json Json;
	if CheckRegistAccountMatchVerifyCode(account,verifyCode)==false{
		json.AppendBool("status",false)
		json.AppendString("msg","验证码错误")
	}else if CheckIsUserExist(account) == false {
		//
		err:=RegistAccount(account,userID,passwd)
		if err==nil{
			json.AppendBool("status",true)
			json.AppendString("msg","regist success")
			Debug("用户:"+r.Header.Get("account")+"注册成功!")
		}else{
			json.AppendBool("status",false)
			json.AppendString("msg",err.Error())
		}
	}else{
		//否则用户存在
		json.AppendBool("status",false)
		json.AppendString("msg","user exist")
	}
	resp:=json.Get()
	_,err=w.Write([]byte(resp))
	return err
}

//注册
func ProcessRegistHanler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessRegistRequest(r,w)
	if err!=nil{
		Debug("regist error!"+err.Error())
	}
}

//生成注册的验证码
func ProcessRegistVerityCodeHandler(cookie string,w http.ResponseWriter,r*http.Request){
	account:=r.Header.Get("account")
	Debug("用户:"+account+"尝试获取验证码")
	var js Json
	//如果之前生成验证码距今时间间隔小于1分钟
	preTime,exist:=RegistVerifyCodeMap.Load(account)
	now:=time.Now().Unix()
	if err:=CheckAccountLegal(account);err!=nil{
		js.AppendBool("status",false)
		js.AppendString("msg",err.Error())
	}else if exist&&(now-preTime.(int64))<60{
		js.AppendBool("status",false)
		js.AppendString("msg","请勿重复发送验证码!")
	}else{
		code:=GenerateRegistVerifyCode(account)
		//将数据插入数据库
		err:=AddAccountRegistCodeToTable(account,code)
		if err!=nil{
			js.AppendBool("status",false)
			js.AppendString("msg","请重试!")
			Debug(err.Error())
		}else{
			RegistVerifyCodeMap.Store(account,now)
			var email EmailInfo
			email.email=account
			email.subject="一生万物注册验证码"
			email.content=fmt.Sprintf("【一生万物】您的 QQ 邮箱验证码为:%s，5 分钟内有效，感谢您使用一生万物相关服务。",code)
			email.deal=nil
			SendEmail(email)
			js.AppendBool("status",true)
			js.AppendString("msg","验证码已发送")
			Debug("验证码发送成功!")
		}
	}

	_,err:=w.Write([]byte(js.Get()))
	if err!=nil{
		Debug("验证码发送出错!")
	}
}
