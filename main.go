package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wsndshx/Paimon/group"
)

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

func main() {
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
	app.Run(":5800")
}
