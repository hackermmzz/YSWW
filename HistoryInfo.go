package main


import(
	"net/http"
	"strconv"
	"errors"
)
//历史记录
func ProcessHistoryInfoHandler(cookie string,w http.ResponseWriter,r*http.Request){
	err := ResponseHistoryInfo(cookie, w)
	if err != nil {
		Debug(err.Error())
		return
	}
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
			infoAll.AppendJson("KouTu"+"_"+strconv.Itoa(id),js)
		}
	}
	_, err:= w.Write([]byte(infoAll.Get()))
	return err
}