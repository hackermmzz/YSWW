package main
import(
	"sync"
)
//发送的邮件需要的信息
type EmailInfo struct{
	email 	string	//邮箱
	subject	string 	//主题
	content	string 	//内容
	deal	func(string)	//处理函数
}
//
var (
	EmailSockChannel  *SocketConn	//通道
	EmailInfoSendLock	sync.RWMutex
	EmailInfoSendQueue	[]EmailInfo//队列 
)

func SendEmail(email EmailInfo){
	EmailInfoSendLock.Lock()
	defer EmailInfoSendLock.Unlock()
	//
	EmailInfoSendQueue=append(EmailInfoSendQueue,email)
}

//初始化邮件发送协程
func EmailSendInit()error{
	EmailInfoSendQueue=make([]EmailInfo,0)
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
	return nil
}