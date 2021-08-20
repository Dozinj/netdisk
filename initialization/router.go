package initialization
import (
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"netdisk/api"
	_ "netdisk/docs"
	"netdisk/global"
	"netdisk/middleware"
)

// 初始化总路由
func Routers() *gin.Engine {
	var Router = gin.Default()
	gin.SetMode(gin.DebugMode)
	//自动生成接口文档---地址 index.html
	Router.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	global.SugaredLogger.Info("register swagger handler")

	//非权限路由
	PublicGroup := Router.Group("")
	{
		PublicGroup.POST("user/register",api.Register)
		PublicGroup.POST("user/login",api.Login)
	}

	//用户鉴权路由
	PrivateGroup := Router.Group("")
	PrivateGroup.Use(middleware.JwtAuth())
	{
		PrivateGroup.POST("file/video", api.UploadVideo)         //上传视频
		PrivateGroup.POST("file/image", api.UploadImage)         //上传图片
		PrivateGroup.POST("share/qrcode", api.GenerateQrcode)     //二维码分享文件
		PrivateGroup.GET("share/qrcode", api.AnalyzeQrcode)     //二维码解析
		PrivateGroup.POST("share/link", api.GenerateSharingLink) //加密链接分享文件
		PrivateGroup.POST("s/:link", api.AnalyzeShareLink)       //分享链接解析
		PrivateGroup.GET("disk/save", api.Save)                  //保存到网盘
		PrivateGroup.GET("local/save", api.Download)             //附件下载
		PrivateGroup.POST("list/video", api.ListVideos)
		PrivateGroup.POST("list/image", api.ListImages)
		PrivateGroup.PUT("file/video/path", api.ChangeVideoPath)
		PrivateGroup.PUT("/file/image/path", api.ChangeImagePath)
		PrivateGroup.PUT("/file/filename", api.ModifyFileName)
		PrivateGroup.PUT("/file/authority",api.ModifyAuthority)
		//PrivateGroup.DELETE("/file",api.DeleteFile)//删除网盘文件
	}

	global.SugaredLogger.Info("router register success")
	return Router
}


