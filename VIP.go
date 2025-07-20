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
	err:=UpdateUserVIPStatus(user,true)
	if err!=nil{
		Debug(err.Error())
		js.AppendBool("status",false)
		js.AppendString("msg","开通失败")
	}else{
		js.AppendBool("status",true)
		js.AppendString("msg","开通成功")
	}
	//
	_,err=w.Write([]byte(js.Get()))
	if err!=nil{
		Debug(err.Error())
	}
}