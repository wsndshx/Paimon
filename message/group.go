package message

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wsndshx/Paimon/module"
	"github.com/wsndshx/Paimon/utils"
)

var Ai bool

func Handle(message string, num int64, user int64) {
	// 构建消息体
	msg := utils.Reply{
		Message_type: "group",
		Group_id:     num,
	}
	// 设置一个随机数
	rand.Seed(time.Now().Unix())
	// Humiliate_key 消息中包含的派蒙的难听绰号
	Humiliate_key := ""
	// morning 是否为早安
	Morning := false
	// night 是否为晚安
	Night := false
	// Praise 是否在称赞派蒙
	Praise_key := ""
	// Cmd 判断接收到的指令
	Cmd := ""

	// 启用ai时使用api返回值作为判断依据, 否则使用正则表达式进行匹配
	if Ai {
		// 匹配被At到的行为: 作为用户语句中的派蒙使用
		at := false
		{
			at_str := regexp.MustCompile(`\[CQ:at,qq=3381113848\]`)
			if at_str.MatchString(message) {
				at = true
			}
		}

		// 获取分析结果
		data := utils.Analysis(message)
		log.Printf("解析结果:\n意图: %s\n包含的实体: %v\n包含的特征: %v\n", utils.Intents[data.Intents], data.Entities, data.Traits)
		if data.Entities["paimon"] == "" && at {
			data.Entities["paimon"] = "at"
		}
		// 判断用户在说什么
		if data.Entities["paimon"] != "" || at {
			switch data.Intents {
			case 0:
				// 说个什么东西搪塞过去
				content := []string{
					"欸...? 旅行者说的东西好深奥派蒙听不懂......",
					"旅行者你在说什么呀?",
					"欸?",
					"听不懂呢.......",
				}
				msg.Message = content[rand.Intn(len(content))]
				msg.Reply()
				return
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
				log.Println("检测到特殊触发" + utils.Intents[data.Intents] + ", 开始进行关键字匹配...")
				// 获取触发内容
				log.Println("触发关键字为: ", data.Entities["paimon"])
				Humiliate_key = data.Entities["paimon"]

			case 4:
				// 检测到称赞
				log.Println("检测到特殊触发" + utils.Intents[data.Intents] + ", 开始进行关键字匹配...")
				log.Println("触发关键字为: ", data.Entities["paimon"])
				Praise_key = data.Entities["paimon"]
			case 6:
				// 唤醒派蒙
				content := []string{
					"派蒙在哦",
					"我在~",
					"(探头)",
					"旅行者你在叫我吗?",
				}
				msg.Message = content[rand.Intn(len(content))]
				msg.Reply()
				return
			case 7:
				// 触发抽卡

				var typer string  // 池子的类型
				var pooler string // 选择的池子
				var times string  // 祈愿的次数

				if pooler = data.Entities["role"]; pooler != "" {
					typer = "Role"
				} else if pooler = data.Entities["weapon"]; pooler != "" {
					typer = "Weapon"
				} else {
					pooler = "null"
					typer = "Resident"
				}

				if times = data.Entities["wit$number"]; times == "" {
					times = "1"
				}

				Cmd = fmt.Sprintf("Wish %s %s %s", typer, pooler, times)
				log.Println("检测到指令 `" + Cmd + "` , 开始执行...")
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
			content := []string{
				"哼哼, 那当然了",
				"嘿嘿嘿",
			}
			msg.Message = content[rand.Intn(len(content))]
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

	// 指令解析
	{
		if Cmd != "" {
			arr := strings.Fields(Cmd)

			// 解析一级指令
		OutCMD:
			switch arr[0] {
			case "Wish":
				msg.Message = "少女祈祷中......"
				msg.Reply()

				// 存储祈愿结果
				var result []string

				// 解析抽取的次数
				var times uint8
				if time, err := strconv.ParseUint(arr[3], 10, 8); err != nil {
					msg.Message = "呜呜呜出错了: " + err.Error()
					break OutCMD
				} else {
					times = uint8(time)
				}

				// 解析二级指令
				switch arr[1] {
				case "Resident":
					// 常驻祈愿
					if data, err := module.Resident(times, uint64(user)); err != nil {
						msg.Message = "呜呜呜出错了: " + err.Error()
						break OutCMD
					} else {
						result = data
					}
				case "Role":
					// 角色祈愿
					if data, err := module.Role(times, uint64(user), arr[2]); err != nil {
						msg.Message = "呜呜呜出错了: " + err.Error()
						break OutCMD
					} else {
						result = data
					}
				case "Weapon":
					// 武器祈愿
				}

				// 当祈愿次数小于等于20时使用纯文本展示,
				// 否则使用网站链接展示
				if times <= 20 {
					msg.Message = fmt.Sprintf("太好了旅行者, 抽到了这些东西呢: \n%v", result)
				} else {
					msg.Message = fmt.Sprintf("太好了旅行者, 抽到了这些东西呢: \n%s\n详细数据: %s", result[0], result[1])
				}
			}
			// 输出指令的执行结果
			msg.Reply()
		}
	}
}
