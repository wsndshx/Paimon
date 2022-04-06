package message

import (
	"fmt"
	"strings"

	"github.com/wsndshx/Paimon/utils"
)

func Private(message string, id int64) {
	msg := utils.Reply{
		Message_type: "private",
		User_id:      int64(id),
	}
	if message[:1] == "/" {
		cmd := strings.Fields(message[1:])
		switch cmd[0] {
		case "cron":
			switch cmd[1] {
			case "list":
				msg.Message = "目前包含的任务为: \n" + Timer.CronList()

			}
		default:
			msg.Message = fmt.Sprintf("指令在位置%d上解析错误: 不存在指令%s", 0, cmd[1])
		}
		msg.Reply()
		return
	}
	// 私聊数据全部扔给ai处理
	data := utils.Analysis(message)

	// 输出
	if data.Intents == 0 {
		msg.Message = fmt.Sprintf("呜呜, 我听不懂你在说什么. 但我猜:%s", data.Traits)
	}
	msg.Message = fmt.Sprintf("我认为你在 %s\n其中包含的实体为: %v\n并且可能具有以下特征:%s", utils.Intents[data.Intents], data.Entities, data.Traits)
	msg.Reply()
}
