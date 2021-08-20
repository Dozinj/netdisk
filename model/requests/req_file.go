package requests

type FC =GQrcode
type GLink struct {
	Fileid         string `json:"fileid" form:"fileid" binding:"required"`
	Category       string `json:"category" form:"category" binding:"required,gt=0"`
	ExtractionCode string `json:"extraction_code" form:"extraction_code"` //提取码
}

type ALink struct {
	Fileid                   string `json:"fileid"`
	Category                 string `json:"category"`
	ShareUser                string `json:"share_user"`
	EncryptionExtractionCode string `json:"encryption_extraction_code"` //加密提取码
}

type  GQrcode struct {
	Fileid   string `json:"fileid" form:"fileid" binding:"required" `
	Category string `json:"category" form:"category" binding:"required" `
}
type  AQrcode struct {
	Fileid   string `json:"fileid" form:"fileid" binding:"required,gt=0" `
	Category string `json:"category" form:"category" binding:"required,gte=1,lte=2" `
	Shareuser string `json:"shareuser" form:"shareuser" binding:"required"`
}

type ChangePath struct {
	Fileid   string `json:"fileid" form:"fileid" binding:"required"`
	FilePath string `json:"file_path" form:"file_path" binding:"required"`
}


type ChangeFilename struct {
	Fileid   string `json:"fileid" form:"fileid" binding:"required,gt=0" `
	Category string `json:"category" form:"category" binding:"required,gte=1,lte=2" `
	Filename string `json:"filename" form:"filename" binding:"required"`
}

type ChangeAuthority struct {
	Fileid   string `json:"fileid" form:"fileid" binding:"required" `
	Category string `json:"category" form:"category" binding:"required" `
	Authority string `json:"authority" form:"authority" binding:"required"`
}


