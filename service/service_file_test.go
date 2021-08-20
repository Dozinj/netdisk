package service

import (
	"fmt"
	"strings"
	"testing"
)

func Test_createVideoPath(t *testing.T) {
	type args struct {
		fileid   string
		filePath string
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: testing.CoverMode(),
			args: args{filePath: "/learn/code/Qimi/"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePaths:=strings.Split(tt.args.filePath,"/")

			paths:=make([]string,0,len(filePaths))
			var path string
			//获取多级目录
			for _,tmppath:=range filePaths{
				if tmppath==""{
					continue
				}
				path=path+"/"+tmppath

				folderPath:=strings.Split(path,"/"+tmppath)[0]
				fmt.Println(folderPath,tmppath)
				paths=append(paths,path)
			}

			t.Log(paths)
		})
	}
}

func Test_createIPath(t *testing.T) {
	type args struct {
		filePath string
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{username: "xzh", filePath: "/a/b/c/d/"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		filepath := tt.args.filePath
		//username := tt.args.username
		//fmt.Println(username)
		filePaths := strings.Split(filepath, "/")
		//fmt.Println(len(filePaths))

		var path string
		//获取多级目录
		for _, tmppath := range filePaths {
			if tmppath == "" {
				path = "/"
				continue
			} else {
				path = path + tmppath + "/"
			}

			//fmt.Println(path)

			//判断该层目录是否存在
			//	fi := new(model.FileImage)
			//	fi.UploadUserName = username
			//	fi.FilePath = path
			//
			//	ok1 := fi.IsExistFileImage()
			//
			//	If := new(model.ImageFolder)
			//	If.Username = username
			//	If.FolderPath = path
			//
			//	ok2 := If.IsExist()
			//

			//新建文件夹
			fmt.Println(path)
			folderPath := strings.Split(path, "/"+tmppath)[0]

			fmt.Println(folderPath+"/")
			fmt.Println(tmppath,"circle")
		}
	}

}
