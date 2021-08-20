package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"

	"netdisk/global"
)

type AliyunOSS interface {
	UploadFile(file *multipart.FileHeader,strMD5 string)(string,error)
	AppendUpload(file *multipart.FileHeader,strMD5 string)(string,error)
	UploadQrcode(qrCode *qrcode.QRCode)(string,error)
}

type aliyunOSS struct{
}


func GetAliyunOss()*aliyunOSS{
	a:=new(aliyunOSS)
	return a
}

var (
	storageType =oss.ObjectStorageClass(oss.StorageStandard)
	readPermission = oss.ObjectACL(oss.ACLPublicRead)
)

var uploadErr =fmt.Errorf("文件上传失败")


//获取文件格式
func GetFormat(filename string)string{
	slice:=strings.Split(filename,".")
	return strings.ToLower(slice[len(slice)-1])
}

// 截取小数位数
func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

//根据文件格式选择bucket
func NewBucket(file *multipart.FileHeader)(*oss.Bucket,string,error) {
	ossConfig := global.Config.Section("oss")
	client, err := oss.New(ossConfig.Key("end_point").String(), ossConfig.Key("app_key").String(), ossConfig.Key("app_secret").String())
	if err != nil {
		global.SugaredLogger.Error("oss new client failed:", zap.Any("err:", err))
		return nil, "", uploadErr
	}

	//根据文件格式获取对应的Bucket
	var bucketName, bucketUrl string
	if CheckImageType(GetFormat(file.Filename)) {
		bucketName, bucketUrl = ossConfig.Key("images_bucket").String(), ossConfig.Key("images_url").String()
	} else {
		bucketName, bucketUrl = ossConfig.Key("videos_bucket").String(), ossConfig.Key("videos_url").String()
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		global.SugaredLogger.Error("oss client open bucket failed:", zap.Any("err:", err))
		return nil, "", uploadErr
	}
	return bucket, bucketUrl, nil
}

//设置文件元信息
func CreateOptions(strMD5 string)[]oss.Option {
	//文件预览行为时需要绑定的自定义域名进行访问
	expires := time.Date(2049, time.January, 10, 23, 0, 0, 0, time.UTC)
	options := []oss.Option{
		oss.Expires(expires),             // 设置缓存过期时间，
		oss.ObjectACL(oss.ACLPublicRead), //公共读
		oss.ContentLanguage("zh-CN"),     //Object使用简体中文编写
		//oss.ContentDisposition("inline"), //直接在浏览器中打开Object。
		//oss.ContentDisposition("attachment"),//保存到本地
		oss.CacheControl("no-cache"), //现向oss验证,缓存可用时直接访问缓存，缓存不可用时重新向OSS请求。
		oss.ContentMD5(strMD5),//文件MD5效验，防止文件被篡改
	}
	return options
}

//返回资源的绝对路径，相对路径，错误信息
func (a *aliyunOSS) UploadFile(file *multipart.FileHeader,strMD5 string)(string,error) {
	bucket,bucketUrl,err:=NewBucket(file)
	if err != nil {
		return "",err
	}

	//读取上传文件
	f, err := file.Open()
	if err != nil {
		global.SugaredLogger.Error("upload_file open failed:", zap.Any("err:", err))
		return "", uploadErr
	}

	//上传bucket中路径保证唯一性
	ObjectName := filepath.Join(time.Now().Format("2006-01-02")) + "/" + file.Filename

	//小文件采用普通方式上传
	ossOptions:=CreateOptions(strMD5)
	if bucket.BucketName==global.Config.Section("oss").Key("images_bucket").String(){
		//使用的是上传图片的桶，则资源打开方式为直接浏览
		ossOptions=append(ossOptions,oss.ContentDisposition("inline"))//直接在浏览器中打开Object
	}else{
		ossOptions=append(ossOptions,oss.ContentDisposition("attachment"))//保存到本地
	}


	err = bucket.PutObject(ObjectName, f, ossOptions...)
	if err != nil {
		global.SugaredLogger.Error("bucket putObject failed:", zap.Any("err:", err))
		return "", uploadErr
	}


	fileUrl:=bucketUrl + "/" + ObjectName
	return fileUrl,nil
}


//追加上传 实现断点续传
func  (a *aliyunOSS) AppendUpload(file *multipart.FileHeader,strMD5 string)(string,error) {
	bucket, bucketUrl, err := NewBucket(file)
	if err != nil {
		return "", err
	}

	fd, err := file.Open()
	if err != nil {
		global.SugaredLogger.Error("open file failed", zap.Any("err:", err))
		return "", uploadErr
	}
	var nextPos, total int64 = 0, 0
	ObjectName := filepath.Join(time.Now().Format("2006-01-02")) + "/" + file.Filename

	//判断文件是否存在
	isExist, err := bucket.IsObjectExist(ObjectName)
	if err != nil {
		global.SugaredLogger.Error(err)
	}
	if isExist {
		//获取文件http头部信息
		props, err := bucket.GetObjectDetailedMeta(ObjectName)
		if err == nil { //继续之前已上传过的断点上传
			nextPos, err = strconv.ParseInt(props.Get("X-Oss-Next-Append-Position"), 10, 64)
		}
	}
	for {
		p := make([]byte, 1<<20)
		n, err := fd.Read(p)
		if err == io.EOF {
			global.SugaredLogger.Info(file.Filename + "数据读取完毕")
			break
		}
		//追加上传 + md5效验
		m := md5.New()
		m.Write(p)
		nextPos, err = bucket.AppendObject(ObjectName, bytes.NewReader(p),
			nextPos, readPermission, storageType,
			oss.ContentMD5(base64.StdEncoding.EncodeToString(m.Sum(nil))))
		if err != nil {
			global.SugaredLogger.Error("AppendObject failed:", zap.Any("err:", err))
			return "", fmt.Errorf("上传文件中断")
		}
		total += int64(n)
	}

	ossOptions := CreateOptions(strMD5)
	ossOptions = append(ossOptions, oss.ContentDisposition("attachment"))
	err = bucket.SetObjectMeta(ObjectName, CreateOptions(strMD5)...)
	if err != nil {
		global.SugaredLogger.Error("修改文件元信息失败", zap.Any("err:", err))
	}
	fileUrl := bucketUrl + "/" + ObjectName
	return fileUrl, nil
}


func (a *aliyunOSS) UploadQrcode(qrCode *qrcode.QRCode)(string,error) {
	buffer := new(bytes.Buffer)
	err := qrCode.Write(256, buffer)
	if err != nil {
		global.SugaredLogger.Error("Buffer write  failed", zap.Any("err:", err))
		return "", fmt.Errorf("二维码生成失败")
	}

	ossConfig := global.Config.Section("oss")
	client, err := oss.New(ossConfig.Key("end_point").String(), ossConfig.Key("app_key").String(), ossConfig.Key("app_secret").String())
	if err != nil {
		global.SugaredLogger.Error("oss new client failed:", zap.Any("err:", err))
		return "", uploadErr
	}

	//根据文件格式获取对应的Bucket
	var bucketName, bucketUrl string
	bucketName, bucketUrl = ossConfig.Key("images_bucket").String(), ossConfig.Key("images_url").String()
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		global.SugaredLogger.Error("oss client open bucket failed:", zap.Any("err:", err))
		return "", uploadErr
	}

	ObjectName := filepath.Join("share-qrcode") + "/"+ strconv.FormatInt(time.Now().Unix(), 10)+".jpg"
	ossOptions:=CreateOptions("")
	ossOptions=append(ossOptions,oss.ContentDisposition("inline"))
	err=bucket.PutObject(ObjectName,buffer,ossOptions...)
	if err!=nil{
		global.SugaredLogger.Error("bucket putObject failed:", zap.Any("err:", err))
		return "", uploadErr
	}

	return bucketUrl+"/"+ObjectName,nil
}

