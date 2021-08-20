

### 实现功能

* 用户登录注册
* 上传文件(暂时只支持视频和文件)
* 秒传
* 文件修改路径，重命名
* 加密链接分享
* 二维码分享
* 下载文件
* 权限管理




## 项目地址 
* http://xuzhihao.cn


## 接口文档
[文档地址](http://xuzhihao.cn/swagger/index.html)  



## 架构设计

```lua
├── api
|   ├── api_user.go  用户相关
|   └── api_file_share.go  文件分享相关
|   └─  api_file.go  文件相关
├── config -- 配置文件 
|   ├── config.ini
├── core 
|   ├── server.go  httpserver 
├── docs -- swagger文档
|   ├── docs.go 
|   ├── swagger.json -- json
|   └── swagger.yaml -- yaml  
├── global -- global
|   |──global.go 全局变量
├── initialize 启动前初始化
|   |──gorm.go
|   |──ini.go   配置
|   |__redis.go 
|   |__router.go  路由组
|   |__validator.go 请求效验库
|   |___zap.go 日志库
├── middleware -- 中间件
|   |___jwt.go
|   |___cors.go
├── model   包含数据库相关操作
│   ├── request  -- 请求参数绑定
|   |   ├── req_file.go
|   |   └── req_user.go  
|   ├── response  -- 返回参数
|   |   ├── resp_file.go 
|   |   ├── resp_user.go
|   |   └── response.go 返回数据格式
|   |____image.go    图片文件
|   |____image_folder.go  图片文件夹
|   |____image_share.go 文件分享相关
|   |____jwt.go
|   |_____user.go
|   |____video.go  视频
|   |____video_folder.go
|   |_____video_share.go
|___service  服务层
|   |____service_file.go
|   |____service_file_share.go
|   |____service_file_test.go
|   |____service_user.go
├── utils 阿里云oss调用
├── Dockerfile  -- docker配置
└── main.go  
```

