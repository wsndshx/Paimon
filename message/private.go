package message

import (
	"fmt"

	"github.com/wsndshx/Paimon/utils"
)

func Private(message string, id int64) {
	// 私聊数据全部扔给ai处理
	data := utils.Analysis(message)
	msg := utils.Reply{
		Message_type: "private",
		User_id:      int64(id),
	}
	// 输出
	if data.Intents == 0 {
		msg.Message = fmt.Sprintf("呜呜, 我听不懂你在说什么. 但我猜:%s", data.Traits)
	}
	msg.Message = fmt.Sprintf("我认为你在 %s, 并且可能具有以下特征:%s", utils.Intents[data.Intents], data.Traits)
	msg.Reply()
}
