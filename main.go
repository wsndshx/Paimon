package main

import (
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wsndshx/Paimon/group"
	"github.com/wsndshx/Paimon/utils"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Core Core `yaml:"core"`
}
type Core struct {
	Http_post   string `yaml:"http_post"`
	Cqhttp_host string `yaml:"cqhttp_host"`
}

type Message struct {
	Post_type    string `json:"post_type"`
	Message_type string `json:"message_type"`
	Sub_type     string `json:"sub_type"`
	Temp_source  int    `json:"temp_source"`
	Message_id   int32  `json:"message_id"`
	User_id      int64  `json:"user_id"`
	Group_id     int64  `json:"group_id"`
	Raw_message  string `json:"raw_message"`
	Font         int32  `json:"font"`
}

func init() {
	log.SetPrefix("[Core]")
	log.SetFlags(0)
}

func (conf *Conf) getConf() *Conf {
	yamlFile, err := ioutil.ReadFile("data/conf.yaml")
	if err != nil {
		log.Panicln(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Panicln(err.Error())
	}
	return conf
}

func main() {
	// 读取配置文件
	var conf Conf
	conf.getConf()
	utils.Host = conf.Core.Cqhttp_host
	// 监听post请求
	app := gin.Default()
	app.POST("/", func(c *gin.Context) {
		// 获取接收的消息
		message := Message{}
		c.BindJSON(&message)
		// 分理消息
		switch message.Message_type {
		case "private":
			// 这里是私聊消息
		case "group":
			// 这里是群聊消息
			group.Handle(message.Raw_message, message.Group_id)
		}
	})
	app.Run(":" + conf.Core.Http_post)
}
