package group

import (
	"regexp"
)

func Handle(message string, num int64) {
	// 匹配早安消息
	morning := regexp.MustCompile(`^(早安|早上好|早鸭) | (早安|早上好|早鸭)$`)
	// 匹配晚安消息
	night := regexp.MustCompile(`^(晚安|我睡了|睡) | (晚安|我睡了)$`)
	if morning.MatchString(message) {
		greeting(1, num)
	} else if night.MatchString(message) {
		greeting(2, num)
	}
}
