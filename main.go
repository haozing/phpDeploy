package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"tarsphpdeploy/config"
	"tarsphpdeploy/tar"
	"tarsphpdeploy/tars"
	"tarsphpdeploy/util"
	"time"
)

func main()  {

	//是否是清理 tgz
	//util.ClearTgz()
	//读取配置文件
	if !config.Config.ReadFile() {
		return
	}
	//将文件过滤后复制到temp目录
	copyFileList(config.Config.RootPath)

	//打包压缩文件
	//生成随机文件名
	paths := []string{
		config.Config.TempDepPath,
	}

	//制作名字
	targzName := config.Config.Service + time.Now().Format("_20060102150405")+".tar.gz"
	RootTargzName := filepath.Join(config.Config.RootPath,targzName)
	err := tar.Tarinate(paths, RootTargzName)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	//删除临时文件
	os.RemoveAll(config.Config.TempTimePath)
	//上传
	//判断配置
	if config.Config.TarsUrl == "" || config.Config.Token == "" {
		fmt.Fprintln(os.Stderr, "请配置TarsUrl和Token")
		return
	}
	Tars := tars.Tars{
		Url:   config.Config.TarsUrl,
		Token: config.Config.Token,
	}
	//上传
	fmt.Printf("上传至:%v.%v\n", config.Config.App, config.Config.Service)
	Tars.Upload(RootTargzName)
	//获取服务状态
	Tars.ServerList(config.Config.App, config.Config.Service)
	//删除压缩文件
	_ = os.Remove(RootTargzName)

}
func copyFileList(path string) {

	//先创建一个临时
	if err := os.MkdirAll(config.Config.TempPath, os.ModePerm); err != nil {
		fmt.Println("临时目录创建失败", err)
		return
	}
	fs,_:= ioutil.ReadDir(path)
	for _,file:=range fs{

		if file.IsDir(){
			if inArray(file.Name()+"/",config.Config.Ignore) {
				continue
			}
			util.CopyPath(filepath.Join(path,file.Name()),config.Config.TempPath)
		}else{
			if inArray(file.Name(),config.Config.Ignore) {
				continue
			}
			util.CopyFile(filepath.Join(path,file.Name()),filepath.Join(config.Config.TempPath,file.Name()))
		}
	}
}
func inArray(target string, str_array []string) bool {
	for _, element := range str_array{
		if target == element{
			return true
		}
	}
	return false
}