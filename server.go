package main

import (
	"crypto/tls"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"io"
	"encoding/json"
)

// ///////////////////////////////////////
var StringMap = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var RequestHandler=make(map[string]func(string,http.ResponseWriter,*http.Request))
var allowedOrigins=make(map[string]bool)
var IsHTTPS=true
var ServerPort=2222
// ///////////////////////////////////////
func InitServer() {
	//初始化所有可以访问服务器的域
	InitAllAllowOrigins()
	//
	port:=":"+strconv.Itoa(ServerPort)
	//设置CPU个数
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)
	//
	mux := http.NewServeMux()
	mux.HandleFunc("/",handler)
	handler_ := corsMiddleware(mux)
	//http.HandleFunc("/", handler)
	if !IsHTTPS{
		Debug("server is running!")
		err := http.ListenAndServe(port, handler_)
		if err != nil {
			log.Fatal(err)
		}
	}else {
		cert, err := tls.LoadX509KeyPair("SSL/cert.pem", "SSL/key.pem")
		if err != nil {
			log.Fatal(err)
		}

		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		server := &http.Server{
			Addr:      port,
			Handler:	handler_,
			TLSConfig: config,
		}
		//
		Debug("server is running!")
		//等待GPU服务器连接
		//
		err = server.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatal(err)
		}
	}
}
func CheckOriginAllowed(origin  string)bool{
	legal,exist:=allowedOrigins[origin]
	/////////////////////////////////////////
	return true;//调试模式默认所有域名全部可以访问我的服务器
	/////////////////////////////////////////
	return exist&&legal
}
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//获取请求头里面的域
		origin := r.Header.Get("Origin")
		legal_origin:="https://www.adhn.asia"
		if CheckOriginAllowed(origin){
			legal_origin=origin
		}
		//
        w.Header().Set("Access-Control-Allow-Origin", legal_origin)
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With,verifyCode, account, type,password,type2,new_passwd")
        w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // 如果需要凭证

		 if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}



func handler(w http.ResponseWriter, r *http.Request) {
	//
	requestType := r.Header.Get("type")
	cookie:=r.Header.Get("cookie")
	exist,user:=CheckCookieLegal(cookie)
	//打印访问记录(对内部cookie带来的访问不提示)
	if exist&&!CheckCookieUsedPrivate(cookie){
		Debug("用户: "+user+" 来访")
	}
	//如果是获取文件资源来的
	if r.URL.Path!="/"{
		requestType="DownLoadFile"
	}else if requestType == "UploadPerson" || requestType == "UploadClothes" {
		requestType="XNCY"
	}
	//更具调用表执行调用函数
	handler_,exist:=RequestHandler[requestType]
	if exist{
		handler_(cookie,w,r)
	}else{
		Debug("request操作:"+requestType+"不存在!")
	}
}	



func RandomFileName() string {
	res := make([]byte, DataBaseFileNameMaxLength)
	for i := 0; i < DataBaseFileNameMaxLength; i++ {
		randIdx := rand.Intn(len(StringMap))
		res[i] = byte(StringMap[randIdx])
	}
	return string(res)
}
func RandomFileNameWithSuffix(suffix string)string{
	l:=len(suffix)
	res := make([]byte, DataBaseFileNameMaxLength-l)

	for i := 0; i < DataBaseFileNameMaxLength-l; i++ {
		randIdx := rand.Intn(len(StringMap))
		res[i] = byte(StringMap[randIdx])
	}
	return string(res)+suffix
}
//回复一个状态值
func RespondStatus(w http.ResponseWriter,status bool)error{
	var js Json
	js.AppendBool("status",status)
	_,err:=w.Write([]byte(js.Get()))
	return err
}
//回复一个收到请求的确认包
func RespondACK(w http.ResponseWriter)error{
	return RespondStatus(w,true)
}
//回复一个异常的确认包
func RespondNCK(w http.ResponseWriter)error{
	return RespondStatus(w,false)
}
//将请求转换成json格式
func ProcessJsonRequestBody(r *http.Request)(map[string]interface{},error){
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil,err
	}
	//关闭请求体
	defer r.Body.Close()
	// 处理请求体（示例：解析为 JSON）
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil,err
	}
	//
	return data,nil
}
//配置处理函数
func ConfigRequestHandler(){
	RequestHandler["DownLoadFile"]=ProcessDownLoadFileHandler;
	RequestHandler["XNCY"]=ProcessXNCYHandler;
	RequestHandler["HistoryInfo"]=ProcessHistoryInfoHandler;
	RequestHandler["WenShengTu"]=ProcessWenShengTuHandler;
	RequestHandler["RLFGH"]=ProcessRlfghHandler;
	RequestHandler["PORTRAIT"]=ProcessPortraitHandler;
	RequestHandler["KOUTU"]=ProcessKouTuHandler;
	RequestHandler["Home"]=ProcessHomeHandler;
	RequestHandler["PlayCountIncrease"]=ProcessPlayCountHandler;
	RequestHandler["PlayCountQuery"]=ProcessPlayCountQueryHanlder;
	RequestHandler["Login"]=ProcessLoginHandler;
	RequestHandler["Regist"]=ProcessRegistHanler;
	RequestHandler["PassWordChange"]=ProcessPassWordChangeHandler;
	RequestHandler["UpLoadFile"]=ProcessUpLoadFileHandler;
	RequestHandler["RegistVerityCode"]=ProcessRegistVerityCodeHandler;
	RequestHandler["VIP"]=ProcessVIPHandler;
	RequestHandler["SendEmail"]=ProcessSendEmailRequestHandler;
	RequestHandler["AIModelTask"]=ProcessAIModelTasksHandler;
	RequestHandler["Tip"]=ProessTipHandler;
	RequestHandler["UserHome"]=ProcessUserHomeHandler;
	RequestHandler["userIDChange"]=ProcessUserIDChangeHandler;
	RequestHandler["avatarChange"]=ProcessAvatarChangeHandler;
	RequestHandler["UserFeedback"]=ProcessUserFeedbackHandler;
	RequestHandler["deleteAccount"]=ProcessDeleteAccountHandler;
}

func InitAllAllowOrigins(){
	allowedOrigins= map[string]bool{
		"http://www.adhn.asia": true,
		"https://www.adhn.asia": true,
	}
}
//