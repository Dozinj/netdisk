package utils

import "sync"

var VideoType []string

var once sync.Once
func getVideoType() *[]string {
	once.Do(func() {
		VideoType=make([]string,0,8)
		VideoType=append(VideoType,"avi","mp4","mpeg","rm","flv","asf","mov","wmv")
	})
	return &VideoType
}

var ImageType []string
func getImageType() *[]string {
	once.Do(func() {
		ImageType=make([]string,0,8)
		ImageType=append(VideoType,"jpg","png","jpeg","gif","bmp","webp","pcx","tif")
	})
	return &ImageType
}



func CheckVideoType(fileType string)bool{
	videoType:=getVideoType()
	for _,v:=range *videoType{
		if v==fileType{
			return true
		}
	}
	return false
}
func CheckImageType(fileType string)bool{
	imageType:=getImageType()
	for _,v:=range *imageType{
		if v==fileType{
			return true
		}
	}
	return false
}

