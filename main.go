package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wsndshx/Paimon/message"
	"github.com/wsndshx/Paimon/utils"
	"gopkg.in/yaml.v2"
)

var timer *utils.Cron

type Conf struct {
	Core   Core   `yaml:"core"`
	Ai     Ai     `yaml:"ai"`
	Notion Notion `yaml:"notion"`
}
type Core struct {
	Http_post   string `yaml:"http_post"`
	Cqhttp_host string `yaml:"cqhttp_host"`
}
type Ai struct {
	Enable bool   `yaml:"enable"`
	Token  string `yaml:"token"`
}
type Notion struct {
	Token            string            `yaml:"token"`
	Wish_result_id   string            `yaml:"wish_result_id"`
	Wish_database_id string            `yaml:"wish_database_id"`
	UserList         map[uint64]string `yaml:"userList"`
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
	var post string
	{
		var conf Conf
		conf.getConf()
		utils.Host = conf.Core.Cqhttp_host
		post = conf.Core.Http_post

		utils.Notion_token = conf.Notion.Token
		utils.Wish_result_id = conf.Notion.Wish_result_id
		utils.Wish_database_id = conf.Notion.Wish_database_id
		for k, v := range conf.Notion.UserList {
			utils.UserList[k] = v
		}

		if conf.Ai.Enable {
			utils.Ai_token = conf.Ai.Token
			utils.Ai_init()
			message.Ai = true
		}
	}

	gin.SetMode(gin.ReleaseMode)
	// gin.DefaultWriter = ioutil.Discard

	// 初始化定时器
	timer = utils.NewCron()
	message.Timer = timer
	if err := timer.Local(); err != nil {
		log.Panic("无法加载定时器数据库: " + err.Error())
	}
	defer timer.Close()

	// 监听post请求
	app := gin.Default()
	app.POST("/", func(c *gin.Context) {
		// 获取接收的消息
		msg := Message{}
		c.BindJSON(&msg)
		text := strings.Replace(msg.Raw_message, "\n", "", -1)
		text = strings.Replace(text, "\r", "", -1)
		// 分理消息
		switch msg.Message_type {
		case "private":
			// 这里是私聊消息

			log.Println("接收到私聊消息: " + text)
			if message.Ai {
				message.Private(text, msg.User_id)
			}
		case "group":
			// 这里是群聊消息
			log.Println("接收到群组消息: " + text)
			message.Handle(text, msg.Group_id, msg.User_id)
		}
	})
	app.Run(":" + post)
}
