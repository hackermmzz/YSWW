package main
import(
	"net/http"
	"os"
	"time"
)

func ProcessUserFeedbackHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//
	legal,user:=CheckCookieLegal(cookie)
	if !legal{
		Debug("cookie不合法")
		return
	}
	//
	js,err:=ProcessJsonRequestBody(r)
	if err!=nil{
		Debug("获取反馈内容失败!")
		return
	}
	content,legal:=js["content"].(string)
	if !legal{
		Debug("获取反馈内容失败!")
		return
	}
	//构建内容
	date:= time.Now().Format("2006-01-02 15:04:05")
	content="时间:"+date+"\n"+"用户:"+user+"\n"+"反馈内容:\n"+content
	//
	filename:=RandomFileName()
	dir:="Source/UserFeedback/"
	feedback,err:=os.Create(dir+filename)
	if err!=nil{
		Debug("创建用户反馈文件失败!")
		return
	}
	//
	_,err=feedback.Write([]byte(content))
	if err!=nil{
		Debug("写入反馈问件失败!")
	}else{
		RespondACK(w)
		Debug("用户:"+user+"发送了一个反馈:"+filename)
	}
}