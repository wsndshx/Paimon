package utils

import (
	"regexp"

	witai "github.com/wit-ai/wit-go/v2"
)

var client *witai.Client
var Ai_token string

// Analysis 语句的分析结果
// Entities 包含的实体
// Traits 词语中包含的特征
// Intents 用户陈述这句话的目的
type AnalysisData struct {
	Intents  uint8
	Traits   []string
	Entities map[string]string
}

func Ai_init() {
	// 初始化连接
	client = witai.NewClient(Ai_token)
	// Use client.SetHTTPClient() to set custom http.Client
}

// Analysis 对用户输入进行分析, 返回语句的分析结果
func Analysis(input string) AnalysisData {
	// 删除无意义的内容
	rm := regexp.MustCompile(`\s{0,1}\[CQ:.*\]\s{0,1}`)
	// 先对内容进行分词
	cws := rm.ReplaceAllString(input, "")
	if cws == "" {
		return AnalysisData{}
	}
	slices := ChineseWS{
		Content: cws,
	}

	// 获取分析结果
	msg, _ := client.Parse(&witai.MessageRequest{
		Query: slices.ChineseWS(),
	})

	data := AnalysisData{}
	// 处理Traits(特征)数据
	if msg.Traits != nil {
		for k, _ := range msg.Traits {
			data.Traits = append(data.Traits, k)
		}
	}
	// 处理Entities(实体)数据
	data.Entities = make(map[string]string)
	if msg.Entities != nil {
		for _, v := range msg.Entities {
			// k 识别到的实体规则名称
			for _, entity := range v {
				// Name 实体的名称
				// Body 识别来源
				data.Entities[entity.Name] = entity.Body
			}
		}
	}
	if len(msg.Intents) == 0 {
		// return fmt.Sprintf("呜呜, 我听不懂你在说什么. 但我猜:%s", Traits)
		data.Intents = 0
		return data
	}
	// return fmt.Sprintf("我认为你在 %s, 并且可能具有以下特征:%s", msg.Intents[0].Name, Traits)
	intents := map[string]uint8{
		"Good_morning": 1,
		"Good_night":   2,
		"Humiliate":    3,
		"Praise":       4,
		"Stating":      5,
		"Awaken":       6,
	}
	data.Intents = intents[msg.Intents[0].Name]
	return data
}
