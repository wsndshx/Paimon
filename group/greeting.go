package group

import "time"

// 问候类型:
// 1 - 早安
// 2 - 晚安
func greeting(Type int8, num int64) {
	t := time.Now().Hour()
	switch {
	case (t < 4):
		switch Type {
		case 1:
			// 回复: 早....欸???
			reply(num, "早....欸???")
		case 2:
			// 是刚刚打完深渊吗(笑), 旅行者晚安~
			reply(num, "是刚刚打完深渊吗(笑), 旅行者晚安~")
		}
	case (t < 7):
		switch Type {
		case 1:
			// 回复: (哈欠)旅行者....早啊......
			reply(num, "(哈欠)旅行者....早啊......")
		case 2:
			// 晚安...欸?
			reply(num, "晚安...欸?")
		}
	case (t < 10):
		switch Type {
		case 1:
			// 回复: 旅行者早鸭~!
			reply(num, "旅行者早鸭~!")
		case 2:
			// 晚安...欸??
			reply(num, "晚安...欸??")
		}
	case (t < 13):
		switch Type {
		case 1:
			// 回复: 早~我发现旅行者自从到了异世后就变懒了呢.......
			reply(num, "早~我发现旅行者自从到了异世后就变懒了呢.......")
		case 2:
			// 晚安...欸???
			reply(num, "晚安...欸???")
		}
	case (t < 19):
		switch Type {
		case 1:
			// 回复: 早~我发现旅行者自从到了异世后就变懒了呢.......
			reply(num, "早~我发现旅行者自从到了异世后就变懒了呢.......")
		case 2:
			// 晚安...?
			reply(num, "晚安...?")
		}
	case (t < 22):
		switch Type {
		case 1:
			// 回复: 旅行者早...欸??
			reply(num, "旅行者早...欸??")
		case 2:
			// 旅行者累了一天了吧.....晚安辣~!
			reply(num, "旅行者累了一天了吧.....晚安辣~!")
		}
	default:
		switch Type {
		case 1:
			// 回复: 旅行者早...欸??!
			reply(num, "旅行者早...欸??!")
		case 2:
			// 好晚了呢...晚安~
			reply(num, "好晚了呢...晚安~")
		}
	}
}
