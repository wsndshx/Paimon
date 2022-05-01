package message

import (
	"fmt"
	"strconv"
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
		errInt := 0
		errStr := ""
		isERR := false
		switch cmd[0] {
		case "cron":
			switch cmd[1] {
			case "list":
				msg.Message = "目前包含的任务为: \n" + Timer.List()
			case "remove":
				errInt++
				if len(cmd) != 3 {
					errInt++
					errStr = fmt.Sprintf("当前指令接受2个参数，却得到了%d个参数", len(cmd)-1)
					break
				}
				// 删除任务
				id, err := strconv.Atoi(cmd[2])
				if err != nil {
					isERR = true
					errStr = "解析任务id时出错: " + err.Error()
					errInt++
				}
				Timer.Remove(id)
			default:
				isERR = true
			}
			errInt++
		default:
			isERR = true
		}
		if isERR {
			if errStr == "" {
				errStr = "不存在指令" + cmd[errInt]
			}
			msg.Message = fmt.Sprintf("指令在位置%d上解析错误: %s", errInt, errStr)
		}
		msg.Reply()
		return
	}
	// 私聊数据全部扔给ai处理
	data, err := utils.Analysis(message)
	if err != nil {
		msg.Message = "呜呜呜, 出错了: " + err.Error()
		msg.Reply()
		return
	}

	// 输出
	if data.Intents == 0 {
		msg.Message = fmt.Sprintf("呜呜, 我听不懂你在说什么. 但我猜:%s", data.Traits)
	}
	msg.Message = fmt.Sprintf("我认为你在 %s\n其中包含的实体为: %v\n并且可能具有以下特征:%s", utils.Intents[data.Intents], data.Entities, data.Traits)
	msg.Reply()
}
