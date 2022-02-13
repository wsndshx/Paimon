package main

import "github.com/wsndshx/Paimon/utils"

func private(message string, id int64) {
	// 私聊数据全部扔给ai处理
	msg := utils.Reply{
		Message_type: "private",
		User_id:      int64(id),
		Message:      utils.Ai(message),
	}
	msg.Reply()
}
