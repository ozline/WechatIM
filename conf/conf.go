package conf

import (
	"fmt"
	"wechat/structs"
	"io/ioutil"

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
	fmt.Println("获取配置文件成功")
}
