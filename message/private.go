package message

import (
	"fmt"
	"strings"

	"github.com/wsndshx/Paimon/module"
	"github.com/wsndshx/Paimon/utils"
)

func Private(message string, id int64) {
	msg := utils.Reply{
		Message_type: "private",
		User_id:      int64(id),
	}
	if message[:1] == "/" {
		cmd := strings.Fields(message[1:])
		module.Cmd(cmd, &msg, uint64(id))
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
