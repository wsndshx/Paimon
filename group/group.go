package group

import (
	"log"
	"regexp"
)

func init() {
	log.SetPrefix("[Group]")
	log.SetFlags(0)
}

func Handle(message string, num int64) {
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
