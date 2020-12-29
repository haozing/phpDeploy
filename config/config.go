package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type config struct {
	App         string            `json:"app"`
	Service     string            `json:"server"`
	Token       string            `json:"token"`
	TarsUrl     string            `json:"tars_url"`
	Ignore     []string 		  `json:"ignore"`
	Path      string              `json:"-"`
	RootPath      string              `json:"-"`
	TempPath      string              `json:"-"`
	TempTimePath      string              `json:"-"`
	TempDepPath      string              `json:"-"`
}

func (c *config) ReadFile() bool {
	//获取配置
	dir,_ := os.Getwd()
	c.Path = filepath.Join(dir,"deployConfig.json")

	file, err := os.Open(c.Path)
	if err != nil {
		fmt.Println("没有找的配置文件")
		return false
	}
	tempConfig := new(config)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(tempConfig)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "解析配置失败:", err)
		return false
	}
	//验证配置是否正确
	if tempConfig.App == "" {
		fmt.Println("app[应用名]未配置")
		return false
	}
	if tempConfig.Service == "" {
		fmt.Println("server[服务名]未配置")
		return false
	}
	c.App = tempConfig.App
	c.Service = tempConfig.Service
	c.Token = tempConfig.Token
	c.TarsUrl = tempConfig.TarsUrl
	c.Ignore = tempConfig.Ignore
	c.RootPath = filepath.Dir(c.Path)
	c.TempTimePath = time.Now().Format("tarsphpdeploytemp_20060102150405")
	c.TempDepPath = filepath.Join(c.RootPath,c.TempTimePath,c.Service)
	c.TempPath = filepath.Join(c.TempDepPath,"src")
	c.Ignore = append(c.Ignore,c.TempTimePath+"/")
	c.Ignore = append(c.Ignore,"phpDeploy")
	c.Ignore = append(c.Ignore,"phpDeploy.exe")

	return true
}


var (
	Config config
)
