package main
import(
	"net/http"
)
//开通会员
func ProcessVIPHandler(cookie string,w http.ResponseWriter,r*http.Request){
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		Debug("cookie不存在")
		return
	}
	//目前不收取任何费用，可以直接开通
	var js Json
	err0:=UpdateUserVIPStatus(user,true)//更新用户会员信息
	err1:=AddVIPInfo(user)//写入会员信息记录表
	if err0!=nil{
		Debug(err0.Error())
		js.AppendBool("status",false)
		js.AppendString("msg","开通失败!")
	}else if err1!=nil{
		js.AppendBool("status",false)
		js.AppendString("msg","请勿重复开通会员!")
	}else{
		js.AppendBool("status",true)
		js.AppendString("msg","开通成功!")
	}
	//
	_,err:=w.Write([]byte(js.Get()))
	if err!=nil{
		Debug(err.Error())
	}
}