package message

import (
	"fmt"
	"log"
	"regexp"

	"github.com/wsndshx/Paimon/utils"
)

var Ai bool
var intents []string

func init() {
	log.SetPrefix("[Group]")
	log.SetFlags(0)
	intents = []string{
		"null",
		"Good_morning",
		"Good_night",
		"Humiliate",
		"Praise",
		"Stating",
		"Awaken",
	}
}

func Handle(message string, num int64) {
	// 构建消息体
	msg := utils.Reply{
		Message_type: "group",
		Group_id:     num,
	}
	// Humiliate_key 消息中包含的派蒙的难听绰号
	Humiliate_key := ""
	// morning 是否为早安
	Morning := false
	// night 是否为晚安
	Night := false
	// Praise 是否在称赞派蒙
	Praise_key := ""

	// 启用ai时使用api返回值作为判断依据, 否则使用正则表达式进行匹配
	if Ai {
		// 匹配被At到的行为: 转Ai, 且不再匹配后续规则
		at := regexp.MustCompile(`\[CQ:at,qq=3381113848\]`)
		if at.MatchString(message) {
			// 获取分析结果
			data := utils.Analysis(message)
			// 输出
			if data.Intents == 0 {
				msg.Message = fmt.Sprintf("呜呜, 我听不懂你在说什么. 但我猜:%s", data.Traits)
			}
			msg.Message = fmt.Sprintf("我认为你在 %s, 并且可能具有以下特征:%s", intents[data.Intents], data.Traits)
			msg.Reply()
			return
		}
		// 获取分析结果
		data := utils.Analysis(message)
		// 判断用户在说什么
		if data.Entities["paimon"] != "" {
			switch data.Intents {
			case 1:
				// 用户输入为早安
				log.Println("消息`" + message + "`为早安问候")
				greeting(1, num)
			case 2:
				// 用户输入为晚安
				log.Println("消息`" + message + "`为晚安问候")
				greeting(2, num)
			case 3:
				// 用户在念派蒙的绰号
				log.Println("检测到特殊触发" + intents[data.Intents] + ", 开始进行关键字匹配...")
				// 获取触发内容
				log.Println("触发关键字为: ", data.Entities["paimon"])
				Humiliate_key = data.Entities["paimon"]
			case 4:
				// 检测到称赞
				log.Println("检测到特殊触发" + intents[data.Intents] + ", 开始进行关键字匹配...")
				log.Println("触发关键字为: ", data.Entities["paimon"])
				Praise_key = data.Entities["paimon"]
			}
		}
	} else {
		// 匹配特殊触发
		anger := regexp.MustCompile(`臭派蒙|应急食品|神奇海鲜|^诶嘿$`)
		Humiliate_key = anger.FindString(message)
		// 匹配早安消息
		morning := regexp.MustCompile(`^早+|蒙早.{0,10}$`)
		Morning = morning.MatchString(message)
		// 匹配晚安消息
		night := regexp.MustCompile(`^晚安|蒙晚安|睡.{0,1}了`)
		Night = night.MatchString(message)
	}

	// 特殊回复相关
	{
		switch Humiliate_key {
		case "":
		case "臭派蒙":
			msg.Message = "喂, 我才不臭呢! (跺脚)"
		case "应急", "食品":
			msg.Message = "我才不是应急食品!"
		case "神奇", "海鲜":
			msg.Message = "喂! 已经从[应急食品]变成[神奇海鲜]了吗?!"
		case "诶嘿":
			msg.Message = "[诶嘿] 是什么意思啊!"
		default:
			if Ai {
				msg.Message = "哼! 我要给你起个难听的绰号!!"
			}
		}
		if msg.Message != "" {
			// 发送特殊触发内容
			msg.Reply()
		}

		if Praise_key != "" {
			msg.Message = "嘿嘿嘿"
			msg.Reply()
		}
	}

	// 早安晚安
	{
		if Morning {
			log.Println("消息`" + message + "`为早安问候")
			greeting(1, num)
		} else if Night {
			log.Println("消息`" + message + "`为晚安问候")
			greeting(2, num)
		}
	}
}
