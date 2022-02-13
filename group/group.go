package group

import (
	"log"
	"regexp"

	"github.com/wsndshx/Paimon/utils"
)

func init() {
	log.SetPrefix("[Group]")
	log.SetFlags(0)
}

func Handle(message string, num int64) {
	{
		// 匹配特殊触发
		msg := utils.Reply{
			Message_type: "group",
			Group_id:     num,
		}
		anger := regexp.MustCompile(`臭派蒙|应急食品|神奇海鲜|^诶嘿$`)
		ao := anger.FindString(message)
		switch ao {
		case "臭派蒙":
			msg.Message = "喂, 我才不臭呢! (跺脚)"
		case "应急食品":
			msg.Message = "我才不是应急食品!"
		case "神奇海鲜":
			msg.Message = "喂! 已经从[应急食品]变成[神奇海鲜]了吗?!"
		case "诶嘿":
			msg.Message = "[诶嘿] 是什么意思啊!"
		}
		// 发送特殊触发内容
		msg.Reply()
	}
	// 匹配早安消息
	morning := regexp.MustCompile(`^早+|蒙早.{0,10}$`)
	// 匹配晚安消息
	night := regexp.MustCompile(`^晚安|蒙晚安|睡.{0,1}了`)
	if morning.MatchString(message) {
		log.Println("消息`" + message + "`为早安问候")
		greeting(1, num)
	} else if night.MatchString(message) {
		log.Println("消息`" + message + "`为晚安问候")
		greeting(2, num)
	}
}
