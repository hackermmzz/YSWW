package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"unicode"
	"regexp"
)

// ///////////////////////////////////////
var StringMap = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var XNCYDataMap 	sync.Map
var RnfghDataMap	sync.Map
var RequestHandler=make(map[string]func(string,http.ResponseWriter,*http.Request))
var allowedOrigins=make(map[string]bool)
var IsHTTPS=true
// ///////////////////////////////////////
func InitServer() {
	//初始化所有可以访问服务器的域
	InitAllAllowOrigins()
	//
	port:=":14332"
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
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, account, type,password,type2,new_passwd")
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
	requestType := r.Header.Get("type")
	cookie:=r.Header.Get("cookie")
	exist,user:=CheckCookieLegal(cookie)
	//打印访问记录
	if exist{
		Debug("用户: "+user+" 来访")
	}
	//如果是获取文件资源来的
	if r.URL.Path!="/"{
		requestType="UpLoadFile"
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
func ProcessPortraitRequest(cookie string,r*http.Request,w http.ResponseWriter)error{
		//检查cookie是否合法
		exist,user:=CheckCookieLegal(cookie)
		if !exist{
			return errors.New("cookie不存在")
		}
		//
		face:=DownLoadFile(r)
		if face==""{
			return errors.New("文件下载失败")
		}
		arguments:=make(map[string]interface{})
		arguments["face"]=face
		go ProcessAIModel("Portrait",user,arguments)
		return nil;
}
func ProcessWSTRequest(cookie string,r*http.Request,w http.ResponseWriter)error{
		//检查cookie是否合法
		exist,user:=CheckCookieLegal(cookie)
		if !exist{
			return errors.New("cookie不存在")
		}
		//
		wordFileName := ProcessWenShengTu(user, r, w)
		if wordFileName == "" {
			return errors.New("文件下载失败")
		}
		arguments:=make(map[string]interface{})
		arguments["wordFile"]=wordFileName
		go ProcessAIModel("WST",user, arguments)
		return nil
}
func ProcessKouTuRequest(cookie string,r*http.Request,w http.ResponseWriter)error{
		//检查cookie是否合法
		exist,user:=CheckCookieLegal(cookie)
		if !exist{
			return errors.New("cookie不存在")
		}
		//
		image:=DownLoadFile(r)
		if image==""{
			return errors.New("文件下载失败!")
		}
		arguments:=make(map[string]interface{})
		arguments["image"]=image
		go ProcessAIModel("KouTu",user,arguments)
		return nil
}
func ProcessPlayCountQuery(cookie string,r*http.Request,w http.ResponseWriter)error{
	//检查cookie是否合法
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		return errors.New("cookie不存在")
	}
	//
	id:=r.Header.Get("id")
	cnt,err:=QueryPlayCountById(id)
	if err!=nil{
		return errors.New(user +" error query playcount "+err.Error())
	}
	err=ResponsePlayCountQuery(w,cnt)
	if err!=nil{
		return errors.New(user +" error query playcount "+err.Error())
	}
	return nil
}
func ProcessRlfghRequest(cookie string,r* http.Request,w http.ResponseWriter)error{
	//检查cookie是否合法
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		return errors.New("cookie不存在")
	}
	//
	info_,exist:=RnfghDataMap.Load(user)
	if !exist{
		RnfghDataMap.Store(user,make([]string,2))
		info_,exist=RnfghDataMap.Load(user)
	}
	info:=info_.([]string)
	//
	type2:= r.Header.Get("type2")
	if type2=="Face"{
		face:=DownLoadFile(r)
		if face==""{
			return errors.New("下载图片文件失败")
		}
		info[0]=face
	}else if type2=="Description"{
		wordFile:=ProcessRnfghWordFile(user,r,w)
		if wordFile==""{
			return	errors.New("下载描述文件失败")
		}
		info[1]=wordFile
	}
	if (len(info[1])!=0) && (len(info[0])!=0){
		arguments:=make(map[string]interface{})
		arguments["face"]=info[0]
		arguments["wordFile"]=info[1]
		go ProcessAIModel("RLFGH",user,arguments)
		info[1]=""
		info[0]=""
		Debug("人脸风格化正在生成...")
	}
	return nil
}
func ProcessXNCYRequest(r *http.Request,cookie string,requestType string)error{
		//检查cookie是否合法
		exist,user:=CheckCookieLegal(cookie)
		if !exist{
			return errors.New("cookie不存在")
		}
		fileName := DownLoadFile(r)
		if fileName == "" {
			return errors.New("download file error!")
		}
		
		//
		info_,exist:=XNCYDataMap.Load(user)
		if !exist{
			XNCYDataMap.Store(user,make([]string,2))
			info_,exist=XNCYDataMap.Load(user)
		}
		info:=info_.([]string)
		//
		if(requestType=="UploadPerson"){
			info[0]=fileName
		}else if(requestType=="UploadClothes"){
			info[1]=fileName
		}
		if (len(info[1])!=0)&&(len(info[0])!=0){
			arguments:=make(map[string]interface{})
			arguments["person"]=info[0]
			arguments["clothes"]=info[1]
			go ProcessAIModel("XNCY",user,arguments)
			info[1]=""
			info[0]=""
			Debug("虚拟穿衣正在生成...")
		}
		return nil
}
func ProcessUpLoadFile(cookie string,URI string,w http.ResponseWriter)error{
	//
	legal,_:=CheckCookieLegal(cookie)
	legal=true//对于下载文件,目前先不检查cookie是否正确
	if legal{
		err := UpLoadFile(w, URI)
		return err
	}
	return nil
}
func CheckCookieLegal(cookie string)(bool,string){
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
func RegistAccount(account string,password string)error{
	//检测用户和密码是否合法
	err:=CheckAccountLegal(account)
	if err!=nil{
		return err
	}
	err=CheckPassWordLegal(password)
	if err!=nil{
		return err
	}
	//向数据库添加用户
	err=AddUser(account,password,0)
	if err==nil{
		//添加用户到进程
		AllUsers.Store(account,password)
	}
	return err
}
func ProcessRegistRequest(r* http.Request,w http.ResponseWriter)error{
	account := r.Header.Get("account")
	passwd:=r.Header.Get("password")
	var json Json;
	if CheckIsUserExist(account) == false {
		err:=RegistAccount(account,passwd)
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
	_,err:=w.Write([]byte(resp))
	return err
}
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
func ResponsePlayCountQuery(w http.ResponseWriter,cnt int)error{
	cc:=strconv.Itoa(cnt)
	_, err:= w.Write([]byte(cc))
	return err
}
func ResponseHomeInfo(cookie string,w http.ResponseWriter,requestType string) error{
	//
	info:=""
	if requestType=="0"{
		info=HomeInfo
	}else {
		_,exist:=HomeInfoDiv[requestType]
		if exist{
			info=HomeInfoDiv[requestType]
		}
	}
	///////////////////
	if len(info)==0{
		info="{}"
	}
	///////////////////
	_, err:= w.Write([]byte(info))
	return err
}
func ProcessRnfghWordFile(user string,r *http.Request, w http.ResponseWriter) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	wordFileName := RandomFileName()
	localFileName := SourceDirectory + "/" + wordFileName
	out, err := os.Create(localFileName)
	if err != nil {
		fmt.Printf("failed to open the file %s for writing\n", localFileName)
		return ""
	}
	defer out.Close()
	res:=make(map[string]interface{})
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	out.Write([]byte(res["description"].(string)))
	return wordFileName
}
func ProcessWenShengTu(user string, r *http.Request, w http.ResponseWriter) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	wordFileName := RandomFileName()
	localFileName := SourceDirectory + "/" + wordFileName
	out, err := os.Create(localFileName)
	if err != nil {
		fmt.Printf("failed to open the file %s for writing\n", localFileName)
		return ""
	}
	defer out.Close()
	res:=make(map[string]interface{})
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	out.Write([]byte(res["description"].(string)))
	return wordFileName
}

func ResponseHistoryInfo(cookie string, w http.ResponseWriter) error {
	//检查cookie是否合法
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		return errors.New("cookie不存在")
	}
	//
	var infoAll Json
	//获取所有虚拟穿衣的结果
	{
		infos, err := GetXNCYUserInfo(user)
		id:=0
		if err != nil {
			return err
		}
		for _, info := range infos {
			var js Json
			js.AppendString("person",info.person)
			js.AppendString("clothes",info.clothes)
			js.AppendString("generate",info.generate)
			js.AppendString("date",info.date)
			js.AppendBool("status",info.status)
			//
			id+=1
			infoAll.AppendJson("XNCY"+"_"+strconv.Itoa(id),js)
		}
		
	}
	//获取所有文生图结果
	{
		
		infos, err:= GetWSTUserInfo(user)
		id:=0
		if err != nil {
			return err
		}
		for _, info := range infos {
			var js Json
			js.AppendString("description",info.description)
			js.AppendString("generate",info.generate)
			js.AppendString("date",info.date)
			js.AppendBool("status",info.status)
			//
			id+=1
			infoAll.AppendJson("WST"+"_"+strconv.Itoa(id),js)
		}
		
	}
	//获取所有人脸风格化信息
	{
		infos, err:= GetRLFGHUserInfo(user)
		id:=0
		if err != nil {
			return err
		}
		for _, info := range infos {
			var js Json
			js.AppendString("face",info.face)
			js.AppendString("description",info.description)
			js.AppendString("generate",info.generate)
			js.AppendString("date",info.date)
			js.AppendBool("status",info.status)
			//
			id+=1
			infoAll.AppendJson("RLFGH"+"_"+strconv.Itoa(id),js)
		}
	}
	//获取人脸肖像所有信息
	{
		infos,err:=GetPortraitUserInfo(user)
		id:=0
		if err != nil {
			return err
		} 
		for _, info := range infos {
			var js Json
			js.AppendString("person",info.person)
			js.AppendString("generate",info.generate)
			js.AppendString("date",info.date)
			js.AppendBool("status",info.status)
			//
			id+=1
			infoAll.AppendJson("Portrait"+"_"+strconv.Itoa(id),js)
		}
	}
	//获取抠图所有信息
	{
		infos,err:=GetKouTuUserInfo(user)
		id:=0
		if err != nil {
			return err
		} 
		for _, info := range infos {
			var js Json
			js.AppendString("image",info.image)
			js.AppendString("generate",info.generate)
			js.AppendString("date",info.date)
			js.AppendBool("status",info.status)
			//
			id+=1
			infoAll.AppendJson("kouTu"+"_"+strconv.Itoa(id),js)
		}
	}
	_, err:= w.Write([]byte(infoAll.Get()))
	return err
}

func DownLoadFile(r *http.Request) string {
	r.ParseMultipartForm(int64(FileMaxSize))
	mForm := r.MultipartForm
	for k := range mForm.File {
		file, _, err := r.FormFile(k)
		if err != nil {
			fmt.Println("inovke FormFile error:", err)
			return ""
		}
		defer file.Close()
		// store uploaded file into local path
		fileName := RandomFileName()
		localFileName := SourceDirectory + "/" + fileName
		out, err := os.Create(localFileName)
		if err != nil {
			fmt.Printf("failed to open the file %s for writing\n", localFileName)
			return ""
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Printf("copy file err:%s\n", err)
			return ""
		}
		////////////////////////////////////////////////////////////////返回文件名称
		return fileName
		////////////////////////////////////////////////////////////////
	}
	return ""
}
func UpLoadFile(w http.ResponseWriter, fileName string) error {
	filename := SourceDirectory + "/" + fileName
	f, err := os.Open(filename)
	if err != nil {
		return errors.New("failed to open the file:"+filename)
	}
	defer f.Close()
	// 获取文件信息
	info, err := f.Stat()
	if err != nil {
		return errors.New("failed to stat file")
	}
	// 设置响应头，告诉客户端这是一个文件下载
	w.Header().Set("Content - Disposition", "attachment; filename="+info.Name())
	w.Header().Set("Content - Type", "application/octet - stream")
	w.Header().Set("Content - Length", fmt.Sprintf("%d", info.Size()))
	// 将文件内容发送给客户端
	_, err = io.Copy(w, f)
	if err != nil {
		return errors.New(("failed to send file") + err.Error())
	}
	return nil
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
//检查密码是否相同
func CheckPassWordIsSame(current string,expected string)bool{
	return current==expected;
}
//检测用户是否合法
func CheckAccountLegal(account string)error{
	//长度不能超过20并且格式为A@B.C,只能为字母和数字
	if len(account)>20{
		return errors.New("邮箱长度不能超过20!")
	}else{
		re := regexp.MustCompile(`^[a-zA-Z0-9]+@[a-zA-Z0-9]+.com$`)
		match:=re.MatchString(account)
		if !match{
			return errors.New("邮箱格式不正确!")
		}
	}
	return nil
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
//处理下载操作
func ProcessUpLoadFileHandler(cookie string,w http.ResponseWriter,r*http.Request){
	filepath:=r.URL.Path
	err:=ProcessUpLoadFile(cookie,filepath,w)
	if err!=nil{
		Debug("upload file error! "+err.Error())
		//RespondNCK(w)
	}else{
		Debug("upload file "+filepath+" success!")
		//RespondACK(w)
	}
}
//虚拟穿衣处理操作
func ProcessXNCYHandler(cookie string,w http.ResponseWriter,r*http.Request){
	requestType:=r.Header.Get("type")
	err:=ProcessXNCYRequest(r,cookie,requestType)
	if err!=nil{
		Debug("虚拟穿衣出现异常! "+err.Error())
	}
}
//历史记录
func ProcessHistoryInfoHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err := ResponseHistoryInfo(cookie, w)
	if err != nil {
		fmt.Println(err)
		return
	}
}
//文生图
func ProcessWenShengTuHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessWSTRequest(cookie,r,w)
	if err!=nil{
		fmt.Println(err)
	}else{
		Debug("文生图正在生成...")
	}
}
//人脸风格化
func ProcessRlfghHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessRlfghRequest(cookie,r,w)
	if err!=nil{
		fmt.Println(err)
	}
}
//人脸肖像
func ProcessPortraitHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessPortraitRequest(cookie,r,w)
	if err!=nil{
		fmt.Println(err)
	}else{
		Debug("人脸肖像正在生成")
	}
}
//抠图
func ProcessKouTuHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessKouTuRequest(cookie,r,w)
	if err!=nil{
		fmt.Println(err)
	}else{
		Debug("智能抠图正在生成...")
	}
}
//主页
func ProcessHomeHandler(cookie string,w http.ResponseWriter,r*http.Request){
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		Debug("用户不存在!")
		return
	}
	//
	type2:= r.Header.Get("type2")
	err := ResponseHomeInfo(cookie, w,type2)
	if err != nil {
		Debug(user + " error uploading Home info " + err.Error())
		return
	}
}
//浏览量增加
func ProcessPlayCountHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//
	exist,user:=CheckCookieLegal(cookie)
	if !exist{
		Debug("用户不存在!")
		return
	}
	//
	id:=r.Header.Get("id")
	err:=IncreasePlayCountById(id)
	if err!=nil{
		Debug(user +" error increase playcount "+err.Error())
		return
	}
}
//浏览量数量查询
func ProcessPlayCountQueryHanlder(cookie string,w http.ResponseWriter,r*http.Request){
	//
	exist,_:=CheckCookieLegal(cookie)
	if !exist{
		Debug("用户不存在!")
		return
	}
	//
	err:=ProcessPlayCountQuery(cookie,r,w)
	if err!=nil{
		fmt.Println(err)
	}
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
//注册
func ProcessRegistHanler(cookie string,w http.ResponseWriter,r*http.Request){
	err:=ProcessRegistRequest(r,w)
	if err!=nil{
		Debug("regist error!")
	}
}
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
	passwd:=r.Header.Get("new_passwd")
	//如果与之前的密码相同,或者检测密码是否不合法，那么不理睬
	current_passwd,_:=AllUsers.Load(user)
	if CheckPassWordIsSame(current_passwd.(string),passwd){
		js.AppendBool("status",false)
		js.AppendString("msg","新密码与原来的密码相同!")
	}else if err:=CheckPassWordLegal(passwd);err!=nil{
		js.AppendBool("status",true)
		js.AppendString("msg",err.Error())
	}else{
		//修改密码
		err:=ChangePassWord(user,passwd)
		//修改对应记录里面对的密码
		AllUsers.Store(user,passwd)
		//
		if err!=nil{
			fmt.Println(err)
		}else{
			Debug("用户:"+user+"修改密码成功!")
			js.AppendBool("status",true)
			js.AppendString("msg","修改成功!")
		}
	}
	_,err:=w.Write([]byte(js.Get()))
	if err!=nil{
		fmt.Println(err)
	}
}
//配置处理函数
func ConfigRequestHandler(){
	RequestHandler["UpLoadFile"]=ProcessUpLoadFileHandler;
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
}

func InitAllAllowOrigins(){
	allowedOrigins= map[string]bool{
		"http://www.adhn.asia": true,
		"https://www.adhn.asia": true,
	}
}
//