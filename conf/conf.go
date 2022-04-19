package conf

import (
	"fmt"
	"io/ioutil"
	"wechat/global"
	"wechat/structs"

	"gopkg.in/yaml.v2"
)

var Config *structs.Conf

func Init() {
	tmp := new(structs.Conf)
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println("GetConf err", err)
	}
	err = yaml.Unmarshal(yamlFile, tmp)
	if err != nil {
		fmt.Println("GetConf err", err)
	}
	Config = tmp
	global.UnifiedPrintln("获取配置文件成功", nil)
}
