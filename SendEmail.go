package main
import(
	"net/http"
	"strconv"
)
//发送的邮件需要的信息
type EmailInfo struct{
	email 	string	//邮箱
	subject	string 	//主题
	content	string 	//内容
	deal	func(string)	//处理函数(目前暂不使用)
}
//
var (
	EmailSockChannel  *SocketConn	//通道
	EmailInfoSendQueue	SyncQueue[EmailInfo]//队列 
)

func SendEmail(email EmailInfo){
	EmailInfoSendQueue.push(email)
}

//初始化邮件发送协程
func EmailSendInit()error{
	/*(这个方法解耦度太低，目前暂停使用)
	//等待连接
	var err error
	EmailSockChannel,err=MakeSocketConn("邮箱发送",EmailChannelPort)
	if err!=nil{
		return err
	}
	//建立一个协程,处理邮件分发
	go func(){
		for{
			if len(EmailInfoSendQueue)!=0{
				EmailInfoSendLock.Lock()
				//分配出一个任务
				task:=EmailInfoSendQueue[0]
				EmailInfoSendQueue=EmailInfoSendQueue[1:]
				EmailInfoSendLock.Unlock()
				//处理任务
				err:=EmailSockChannel.Write(task.email,task.subject,task.content)
				if err!=nil{
					Debug("发送邮件:"+task.email+"失败!"+err.Error())
					continue
				}
				ret,err:=EmailSockChannel.Read()
				if err!=nil{
					Debug("读取邮件信息异常!"+err.Error())
					continue
				}
				//调用处理函数
				if task.deal!=nil{
					task.deal(ret)
				}
			}
		}
	}()
	*/
	return nil
}
//处理邮箱模块向服务器的请求
func ProcessSendEmailRequestHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//如果不是内部cookie,不做处理
	if !CheckCookieUsedPrivate(cookie){
		Debug("cookie不合法")
		return
	}

	//这里直接发送当前所有的待发送邮件给模块
	EmailInfoSendQueue.hold()
	defer EmailInfoSendQueue.release()
	cnt:=len(EmailInfoSendQueue.queue)
	var resp Json
	for idx,info:=range EmailInfoSendQueue.queue{
		var resp_ Json
		//
		resp_.AppendString("email",info.email)
		resp_.AppendString("subject",info.subject)
		resp_.AppendString("content",info.content)
		//
		resp.AppendJson(strconv.Itoa(idx),resp_)
	}
	//清空队列
	EmailInfoSendQueue.queue=make([]EmailInfo,0)
	_,err:=w.Write([]byte(resp.Get()))
	if err!=nil{
		Debug("发送邮件消息队列失败!")
	}else if cnt!=0{
		Debug("已经处理目前为止所有的邮箱验证!")
	}
}