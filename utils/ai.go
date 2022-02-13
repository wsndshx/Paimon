package utils

import (
	"fmt"
	"regexp"

	witai "github.com/wit-ai/wit-go/v2"
)

var client *witai.Client
var Ai_token string

func Ai_init() {
	// 初始化连接
	client = witai.NewClient(Ai_token)
	// Use client.SetHTTPClient() to set custom http.Client
}

// Ai 对用户输入进行分析
func Ai(input string) string {
	// 删除无意义的内容
	rm := regexp.MustCompile(`\s{0,1}\[CQ:.*\]\s{0,1}`)
	// 先对内容进行分词
	slices := ChineseWS{
		Content: rm.ReplaceAllString(input, ""),
	}

	// Entities 包含的实体
	// Traits 词语中包含的特征
	// Intents 用户陈述这句话的目的
	msg, _ := client.Parse(&witai.MessageRequest{
		Query: slices.ChineseWS(),
	})
	Traits := ""
	for k, v := range msg.Traits {
		Traits = Traits + fmt.Sprintf("\n%s: %s", k, v[0].Value)
	}
	if len(msg.Intents) == 0 {
		return fmt.Sprintf("呜呜, 我听不懂你在说什么. 但我猜:%s", Traits)
	} else {
		return fmt.Sprintf("我认为你在 %s, 并且可能具有以下特征:%s", msg.Intents[0].Name, Traits)
	}
}
