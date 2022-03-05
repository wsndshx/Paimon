package main

import (
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/wsndshx/Paimon/message"
	"github.com/wsndshx/Paimon/utils"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Core Core `yaml:"core"`
	Ai   Ai   `yaml:"ai"`
}
type Core struct {
	Http_post   string `yaml:"http_post"`
	Cqhttp_host string `yaml:"cqhttp_host"`
}
type Ai struct {
	Enable bool   `yaml:"enable"`
	Token  string `yaml:"token"`
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
	if conf.Ai.Enable {
		utils.Ai_token = conf.Ai.Token
		utils.Ai_init()
		message.Ai = true
	}
	gin.SetMode(gin.ReleaseMode)
	// gin.DefaultWriter = ioutil.Discard
	// 这是一个定时器任务, 临时用用
	c := cron.New()

	_, err := c.AddFunc("0 22 * * ?", func() {
		// 这里执行
		log.Println("执行定时任务......")
		// 构建消息体
		msg := utils.Reply{
			Message_type: "group",
			Group_id:     417176143,
			Message:      "现在已经10点了.....旅行者如果还没背单词的话就赶快去背!",
		}
		msg.Reply()
	})
	if err != nil {
		log.Panicf("添加定时器任务失败 : %v\n", err)
	}
	c.Start()
	defer c.Stop()

	// 监听post请求
	app := gin.Default()
	app.POST("/", func(c *gin.Context) {
		// 获取接收的消息
		msg := Message{}
		c.BindJSON(&msg)
		// 分理消息
		switch msg.Message_type {
		case "private":
			// 这里是私聊消息
			log.Println("接收到私聊消息: " + msg.Raw_message)
			if conf.Ai.Enable {
				message.Private(msg.Raw_message, msg.User_id)
			}
		case "group":
			// 这里是群聊消息
			log.Println("接收到群组消息: " + msg.Raw_message)
			message.Handle(msg.Raw_message, msg.Group_id)
		}
	})
	app.Run(":" + conf.Core.Http_post)
}
