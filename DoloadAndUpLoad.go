package main
import(
	"net/http"
	"fmt"
	"os"
	"io"
	"errors"
)
//本地下载文件
func DownLoadFile(r *http.Request) string {
	r.ParseMultipartForm(int64(FileMaxSize))
	//标记是否以原文件名保存在服务器
	type2:=r.Header.Get("type2")
	var originalFileName bool
	if type2=="true"{
		originalFileName=true
	}else{
		originalFileName=false
	}
	
	//
	mForm := r.MultipartForm
	for k := range mForm.File {
		file,header, err := r.FormFile(k)
		if err != nil {
			Debug("inovke FormFile error:"+err.Error())
			return ""
		}
		defer file.Close()
		// store uploaded file into local path
		var fileName string
		if originalFileName{
			fileName=header.Filename
		}else{
			fileName=RandomFileName()
		}
		//
		localFileName := SourceDirectory + "/" + fileName
		out, err := os.Create(localFileName)
		if err != nil {
			Debug("failed to open the file "+localFileName+" for writing")
			return ""
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			Debug("copy file err:"+err.Error())
			return ""
		}
		////////////////////////////////////////////////////////////////返回文件名称
		return fileName
		////////////////////////////////////////////////////////////////
	}
	return ""
}

//处理下载操作
func ProcessDownLoadFileHandler(cookie string,w http.ResponseWriter,r*http.Request){
	filepath:=r.URL.Path
	err:=UpLoadFile(w,filepath)
	if err!=nil{
		Debug("upload file error! "+err.Error())
	}else{
		Debug("upload file "+filepath+" success!")
	}
}

//本地上传文件
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

//处理文件上传到本机的操作
func ProcessUpLoadFileHandler(cookie string,w http.ResponseWriter,r*http.Request){
	//目前该接口只用于接收学校机子下载文件的操作，虽然有漏洞，但是无所谓了
	exist,_:=CheckCookieLegal(cookie)
	if !exist{
		Debug("cookie不存在")
		return
	}
	//
	filename:=DownLoadFile(r)
	Debug("上传文件:"+filename+"成功!")
	return;
}