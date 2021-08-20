package response

//视频上传
type RespUploadVideo struct {
	Status bool                `json:"status" example:"true"`
	Data   RespUploadVideoData `json:"data"`
	Msg    string              `json:"msg" example:"视频上传成功"`
}
type RespUploadVideoData struct {
	VideoId   uint   `json:"视频id" example:"101"`
	VideoName string `json:"视频名称" example:"干饭.mp4"`
	Size      string `json:"内存大小" example:"100MB"`
}


//图片上传
type RespUploadImage struct {
	Status bool                `json:"status" example:"true"`
	Data   RespUploadImageData `json:"data"`
	Msg    string              `json:"msg" example:"图片上传成功"`
}
type RespUploadImageData struct {
	ImageId   uint   `json:"图片id" example:"102"`
	ImageName string `json:"图片名称" example:"恰饭.jpg"`
	Size      string `json:"内存大小" example:"1MB"`
}


//生成二维码
type RespGQrcode struct {
	Status bool    `json:"status" example:"true"`
	Data RespGQrcodeData `json:"data"`
	Msg  string `json:"msg" example:"生成分享文件二维码成功"`
}
type RespGQrcodeData struct {
	Url string `json:"二维码" example:""`
}

//生成加密链接
type RespShareLink struct {
	Status bool            `json:"status" example:"true"`
	Data RespShareLinkData `json:"data"`
	Msg    string          `json:"msg" example:"分享链接生成成功"`
}
type RespShareLinkData struct {
	ShareLink      string `json:"加密链接"`
	ExtractionCode string `json:"提取码"`
}


//分享文件信息
type RespShareFile struct {
	Status bool              `json:"status" example:"true"`
	Data   RespShareFileData `json:"data"`
	Msg    string            `json:"msg" example:"获取分享文件信息成功"`
}
type RespShareFileData struct {
	FileId     string `json:"文件id" example:"103"`
	FileName   string `json:"文件名称" example:"吃饭.jpeg"`
	UploadTime string `json:"文件上传时间" example:"2021-08-20"`
	Size       string `json:"文件大小" example:"10MB"`
	ShareUser  string `json:"分享者"  example:"xzh"`
}



//本地下载
type RespLocalSave struct {
	Status bool            `json:"status" example:"true"`
	Data RespLocalSaveData `json:"data"`
	Msg    string          `json:"msg" example:"获取下载地址成功"`
}
type RespLocalSaveData struct {
	Url string `json:"资源地址"`
}


type RespList struct {
	Status bool           `json:"status" example:"true"`
	Data   []RespListData `json:"data"`
	Msg    string         `json:"msg" example:"获取网盘信息成功"`
}

type RespListFileData struct {
	Fileid    uint   `json:"文件id" example:"100"`
	CreatedAt string `json:"上传时间" example:"2021-08-20 07:00"`
	UpdatedAt string `json:"修改时间" example:"2021-08-20 09:00"`
	FileName  string `json:"文件名" example:"打电动.jpeg"`
	Size      string `json:"大小" example:"2MB"`
	Authority int    `json:"authority" example:"8"`// authority:8.所有人可下载 9.仅获得链接的人可下载 10.仅自己见
}

type RespListFolderData struct {
	CreatedAt string `json:"上传时间" example:"2021-08-20 17:42"`
	UpdatedAt string `json:"修改时间" example:"2021-08-20 18:53"`
	FolderName  string `json:"文件夹名" example:"学习"`
}

type RespListData struct {
	File   []RespListFileData   `json:"文件"`
	Folder []RespListFolderData `json:"文件夹"`
}
