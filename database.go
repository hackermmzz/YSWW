package main

import (
	"database/sql"
	"os"
	"strconv"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"errors"
)

// //////////////////////////////



// ///////////////////////////////////////配置
var (
	ConfigInfo map[string]string
	/////
	SourceProcessProgram string //处理图片数据的程序
	/////
	SourceDirectory string //存储上传过来的图片的文件夹
	FileMaxSize     int    //文件最大大小
	/////
	DataBaseFileNameMaxLength int     //文件名字段最大长度
	db                        *sql.DB //数据库对象
	DataBaseMaxIdleConns      int
	DataBaseName              string
	AllUsers                  sync.Map //所有的用户
	Cookie					  sync.Map //将cookie和用户相关联
	HomeInfo 				  string //主页面显示文件
	HomeInfoDiv					map[string]string//主页面进行分类别
	RequestMaxProcess			  int //可并行处理的最大数量
	fileDebug						bool //是否开启文件记录调试信息
)

////////////////////////////////////////



//
func ConnectToDataBase(username string, password string, database string) {
	var err error
	dsn := username + ":" + password + "@tcp(127.0.0.1:3306)/" + database
	//连接数据集
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		Debug("dsn:"+dsn+" invalid,err:"+err.Error())
		return
	}
	err = db.Ping() //尝试连接数据库
	if err != nil {
		Debug("open "+dsn+" faild,err:"+err.Error())
		return
	}

	db.SetMaxIdleConns(DataBaseMaxIdleConns)
	_, err = db.Exec("use " + DataBaseName)
	if err != nil {
		Debug(err.Error())
	}
	Debug("连接数据库成功~")
}
func InitDatabase() {
	//连接数据库
	ConnectToDataBase(ConfigInfo["UserName"], ConfigInfo["PassWord"], DataBaseName)
	//获取所有的用户
	AllUsers_, err:= GetAllUsers()
	if err != nil {
		Debug("Could not get all users")
		return
	}
	for key,value:=range AllUsers_{
		AllUsers.Store(key,value)
	}
	//获取主页面信息
	HomeInfoDiv=make(map[string]string)
	err=GetHomeInfo()
	if err!=nil{
		Debug("could not get Home info")
		return
	}
}
func GetHomeInfo()error{
	infoMap,infoJson,err := QueryHomeAll()
	if err!=nil{
		return err
	}
	HomeInfo=infoJson.Get()
	for key,js:=range infoMap{
		HomeInfoDiv[key]=js.Get()
	}
	return nil
}
func GetAllUsers() (map[string]string, error) {
	cmd := "select email,password from " + "users;"
	row, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	ret := make(map[string]string)
	var userInfo UserInfo;
	for row.Next() {
		err := row.Scan(&userInfo.email,&userInfo.passwd)
		if err != nil {
			return nil, err
		}
		ret[userInfo.email] = userInfo.passwd
	}
	return ret, nil
}
func AddUser(email string,userID string,passwd string,vip int) error {
	value:=MergeByCommaAndQuo(email,userID,passwd,strconv.Itoa(vip))
	cmd := "insert into " + "users(email,userID,password,vip)" + " values(" +value+");"
	_, err:= db.Exec(cmd)
	if err != nil {
		return err
	}
	return nil
}

func ChangePassWord(user string,passwd string)error{
	cmd:="update users set password="+MarkByQuo(passwd)+"where email="+MarkByQuo(user);
	_,err:=db.Exec(cmd)
	return err
}
func GetXNCYUserInfo(user string) ([]XNCYInfo, error) {
	cmd := "select person,clothes,generate,date,status from XNCY where account="+MarkByQuo(user)
	row, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	ans := make([]XNCYInfo, 0)
	var info XNCYInfo
	for row.Next() {
		err = row.Scan(&info.person, &info.clothes, &info.generate,&info.date,&info.status)
		if err != nil {
			return nil, err
		}
		ans = append(ans, info)
	}
	return ans, nil
}
func GetWSTUserInfo(user string) ([]WSTInfo, error) {
	cmd := "select description,generate,date,status from WST where account=" + MarkByQuo(user)
	row, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	ans := make([]WSTInfo, 0)
	var info WSTInfo
	for row.Next() {
		err = row.Scan(&info.description, &info.generate,&info.date,&info.status)
		if err != nil {
			return nil, err
		}
		ans = append(ans, info)
	}
	return ans, nil
}
func GetRLFGHUserInfo(user string) ([]RLFGHInfo, error) {
	cmd := "select face,description,generate,date,status from RLFGH where account=" + MarkByQuo(user)
	row, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	ans := make([]RLFGHInfo, 0)
	var info RLFGHInfo
	for row.Next() {
		err = row.Scan(&info.face,&info.description,&info.generate,&info.date,&info.status)
		if err != nil {
			return nil, err
		}
		ans = append(ans, info)
	}
	return ans, nil
}
func GetPortraitUserInfo(user string) ([]PortraitInfo, error) {
	cmd := "select person,generate,date,status from Portrait where account=" + MarkByQuo(user)
	row, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	ans := make([]PortraitInfo, 0)
	var info PortraitInfo
	for row.Next() {
		err = row.Scan(&info.person,&info.generate,&info.date,&info.status)
		if err != nil {
			return nil, err
		}
		ans = append(ans, info)
	}
	return ans, nil
}
func GetKouTuUserInfo(user string) ([]KouTuInfo, error) {
	cmd := "select image,generate,date,status from KouTu where account=" + MarkByQuo(user)
	row, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	ans := make([]KouTuInfo, 0)
	var info KouTuInfo
	for row.Next() {
		err = row.Scan(&info.image,&info.generate,&info.date,&info.status)
		if err != nil {
			return nil, err
		}
		ans = append(ans, info)
	}
	return ans, nil
}
func QueryHomeAll() (map[string]*Json,Json,error) {
	dirPath := "./Home/"
	var infoAll Json
	// 读取目录内容
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil,infoAll,err
	}
	// 遍历目录内容
	info:=make(map[string]*Json,0)
	idx:=0
	for _, entry := range entries {
		filePath:=dirPath+entry.Name()
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil,infoAll,err
		}
		parts:=strings.Split(string(content), "\n")
		var infoMap Json
		tp:=""
		for _,line:=range parts{
			line=strings.TrimSpace(line)
			if line==""{
				continue
			}
			elements:=strings.Split(line,":")
			key:=elements[0]
			content:=elements[1]
			if key=="id"{
				id,err:=strconv.Atoi(content)
				if err!=nil{
					return nil,infoAll,err
				}
				infoMap.AppendInt(key,id)
			}else{
				infoMap.AppendString(key,content)
			}

			if key=="type"{
				tp=content
			}
		}

		_,exist:=info[tp]
		if !exist{
			info[tp]=new(Json)
		}
		//给每个记录加一个status,表示资源是否存在
		infoMap.AppendBool("status",true)
		//
		idx+=1
		key:=tp+"_"+strconv.Itoa(idx)
		info[tp].AppendJson(key,infoMap)
		infoAll.AppendJson(key,infoMap)
	}
	return info,infoAll,nil
}
func AddDataToWST(user string, wordFile string, merge string) error {
	cmd := "insert into WST(account,description,generate) values("+MergeByCommaAndQuo(user,wordFile,merge)+")"
	_, err := db.Exec(cmd)
	if err != nil {
		Debug(cmd)
		return err
	}
	return nil
}

func AddDataToRLFGH(user string,face string,wordFile string,merge string) error{
	cmd:="insert into RLFGH(account,face,description,generate) values("+MergeByCommaAndQuo(user,face,wordFile,merge)+");"
	_,err:=db.Exec(cmd)
	if err!=nil{
		Debug(cmd)
		return err
	}
	return nil
}

func AddDataToXNCY(user string, person string, clothes string, merge string) error {
	cmd := "insert into XNCY(account,person,clothes,generate) values("+MergeByCommaAndQuo(user,person,clothes,merge)+");"
	_, err := db.Exec(cmd)
	if err != nil {
		Debug(cmd)
		return err
	}
	return nil
}

func AddDataToPortrait(user string,face string,merge string) error{
	cmd := "insert into Portrait(account,person,generate) values(" +MergeByCommaAndQuo(user,face,merge)+");"
	_, err := db.Exec(cmd)
	if err != nil {
		Debug(cmd)
		return err
	}
	return nil
}


func AddDataToKouTu(user string,image string,merge string) error{
	cmd := "insert into KouTu(account,image,generate) values(" +MergeByCommaAndQuo(user,image,merge)+ ");"
	_, err := db.Exec(cmd)
	if err != nil {
		Debug(cmd)
		return err
	}
	return nil
}
//由于服务器保证文件名一定唯一，所有统一使用这个来更新表格
func UpdateGenerateStatus(generate string,table string)error{
	cmd:="update "+table+" set status=true where generate="+MarkByQuo(generate);
	_,err:=db.Exec(cmd)
	return err
}
func UpdateXNCYStatus(generate string)error{
	return UpdateGenerateStatus(generate,"XNCY")
}

func UpdateWSTStatus(generate string)error{
	return UpdateGenerateStatus(generate,"WST")
}

func UpdatePortraitStatus(generate string)error{
	return  UpdateGenerateStatus(generate,"Portrait")
}

func UpdateRLFGHStatus(generate string)error{
	return  UpdateGenerateStatus(generate,"RLFGH")
}

func UpdateKouTuStatus(generate string)error{
	return  UpdateGenerateStatus(generate,"KouTu")
}
//
func AddNewIdToPlayCountTable(id string) error{
	cmd:="insert into PlayCount values("+id+",1)"//it is not zero because wlh is hero
	_, err := db.Exec(cmd)
	if err!=nil{
		return err
	} 
	return nil
}

func QueryPlayCountById(id string) (int,error){
	cmd:="select count from PlayCount where playCountId="+id
	row,err:= db.Query(cmd)
	defer row.Close()
	if err!=nil{
		return 0,err
	} 
	cnt:=0
	if row.Next(){
		err=row.Scan(&cnt)
	}else{
		err:=AddNewIdToPlayCountTable(id)
		if err!=nil{
			return 0,err
		}
		cnt=1
	}
	return cnt,nil
}

func IncreasePlayCountById(id string) error{
	cmd:="update  PlayCount set count=count+1 where PlayCount="+id
	_, err := db.Exec(cmd)
	if err!=nil{
		return err
	} 
	return nil
}

func QueryAccountRegistCode(account string)([]RegistVerityCodeInfo,error){
	cmd:="select code,date from RegistVerityCode where account="+MarkByQuo(account)+" and status=false "+"and date>=NOW() - INTERVAL 5 minute"
	row,err:=db.Query(cmd)

	if err!=nil{
		return nil,err
	}
	ret:=make([]RegistVerityCodeInfo,0)
	var info RegistVerityCodeInfo
	for row.Next(){
		err=row.Scan(&info.code,&info.date)
		if err!=nil{
			return nil,err
		}
		ret=append(ret,info)
	}
	return ret,nil
}
func UpdateAccountRegistCodeStatus(account string,code string)error{
	cmd:="update RegistVerityCode set status=true where account="+MarkByQuo(account)+" and code="+MarkByQuo(code)
	_,err:=db.Exec(cmd)

	return err;
}

func AddAccountRegistCodeToTable(account string,code string)error{
	cmd:="insert into RegistVerityCode(account,code) values("+MergeByCommaAndQuo(account,code)+")"
	_,err:=db.Exec(cmd)
	return err
}

func UpdateUserVIPStatus(account string,status bool)error{
	status_s:="true"
	if !status{
		status_s="false"
	}
	cmd:="update users set vip="+status_s+" where email="+MarkByQuo(account);
	_,err:=db.Exec(cmd)
	return err
}
//获取对应用户的昵称
func GetUserInfo(account string)(UserInfo,error){
	var info UserInfo
	cmd:="select * from users where email="+MarkByQuo(account)
	row,err:=db.Query(cmd)
	if err!=nil{
		return info,err
	}
	if row.Next(){
		err:=row.Scan(&info.email,&info.userID,&info.passwd,&info.date,&info.vip,&info.avatar)
		if err!=nil{
			return info,err
		}
	}else{
		return info,errors.New("没有对应的用户!")
	}
	return info,nil
}

func AddVIPInfo(user string)error{
	cmd:="insert into vip(account,end_time) values("+MarkByQuo(user)+",DATE_ADD(NOW(), INTERVAL 1 YEAR)"+")"
	_,err:=db.Exec(cmd)
	return err
}

func GetVIPInfo(user string)(VIPInfo,error){
	var info VIPInfo
	cmd:="select * from vip where account="+MarkByQuo(user)
	row,err:=db.Query(cmd)
	if err!=nil{
		return info,err
	}
	if row.Next(){
		err=row.Scan(&info.account,&info.start_time,&info.end_time)
	}else{
		return info,errors.New("没有对应的账号信息!")
	}
	return info,err
}

func UpdateUserID(account string,userID string)error{
	cmd:="update users set userID="+MarkByQuo(userID)+" where email="+MarkByQuo(account);
	_,err:=db.Exec(cmd)
	return err
}

func UpdateUserAvatar(account string,avatarFileName string)error{
	cmd:="update users set avatar="+MarkByQuo(avatarFileName)+" where email="+MarkByQuo(account);
	_,err:=db.Exec(cmd)
	return err
}

func DeleteAccount(account string)error{
	cmd:="delete from users where email="+MarkByQuo(account)
	_,err:=db.Exec(cmd)
	return err
}