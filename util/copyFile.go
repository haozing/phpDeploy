package util

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"tarsphpdeploy/config"
)

//判断文件或目录是否存在
func GetFileInfo(src string) os.FileInfo {
	if fileInfo, e := os.Stat(src); e != nil {
		if os.IsNotExist(e) {
			return nil
		}
		return nil
	} else {
		return fileInfo
	}
}

//拷贝文件
func CopyFile(src, dst string) bool {
	if len(src) == 0 || len(dst) == 0 {
		return false
	}
	srcFile, e := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if e != nil {
		fmt.Println(src+" 创建文件错误", e)
		return false
	}
	defer srcFile.Close()

	//这里要把O_TRUNC 加上，否则会出现新旧文件内容出现重叠现象
	dstFile, e := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if e != nil {
		fmt.Println(dst+" 创建文件错误", e)
		return false
	}
	defer dstFile.Close()
	//fileInfo, e := srcFile.Stat()
	//fileInfo.Size() > 1024
	//byteBuffer := make([]byte, 10)
	if _, e := io.Copy(dstFile, srcFile); e != nil {
		fmt.Println(" 复制文件错误", e)
		return false
	} else {
		return true
	}

}

//拷贝目录
func CopyPath(src, dst string) bool {
	srcFileInfo := GetFileInfo(src)
	if srcFileInfo == nil || !srcFileInfo.IsDir() {
		fmt.Println(src+"不是目录文件")
		return false
	}
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("遍历文件目录错误", err)
			return err
		}

		dstPath := strings.Replace(path,config.Config.RootPath,dst,1)
		if !info.IsDir() {
			if CopyFile(path,dstPath) {
				return nil
			} else {
				return errors.New(path + " 复制文件错误")
			}
		} else {
			if _, err := os.Stat(dstPath); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
						fmt.Println(dstPath+" 创建目录错误", err)
						return err
					} else {
						return nil
					}
				} else {
					fmt.Println("不是正确的目录", err)
					return err
				}
			} else {
				return nil
			}
		}
	})

	if err != nil {
		fmt.Println("CopyPath", err)
		return false
	}
	return true

}