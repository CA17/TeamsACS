package main

import (
	"os"

	"github.com/ca17/teamsacs/config"
	"gopkg.in/yaml.v2"
)

// 初始化一个本地开发配置文件

func main() {
	bs, err := yaml.Marshal(config.DefaultAppConfig)
	if err != nil {
		panic(err)
	}
	os.WriteFile("teamsacs.yml", bs, 777)
}
