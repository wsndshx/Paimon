package utils

import (
	"log"
	"regexp"
	"strconv"

	witai "github.com/wit-ai/wit-go/v2"
)

var client *witai.Client
var Ai_token string
var Intents []string = []string{
	"null",
	"Good_morning",
	"Good_night",
	"Humiliate",
	"Praise",
	"Stating",
	"Awaken",
	"Wish",
	"Reminder",
}

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
func Analysis(input string) (AnalysisData, error) {
	// 删除无意义的内容
	rm := regexp.MustCompile(`\s{0,1}\[CQ:.*\]\s{0,1}`)
	// 先对内容进行分词
	cws := rm.ReplaceAllString(input, "")
	if cws == "" {
		return AnalysisData{}, nil
	}
	slices := ChineseWS{
		Content: cws,
	}

	// 获取分析结果
	msg, err := client.Parse(&witai.MessageRequest{
		Query: slices.ChineseWS(),
	})
	if err != nil {
		log.Fatalln(err)
		return AnalysisData{}, err
	}

	data := AnalysisData{}
	// 处理Traits(特征)数据
	if msg.Traits != nil {
		for k := range msg.Traits {
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
				// Body(Value) 识别来源/结果
				switch entity.Value.(type) {
				case string:
					//string类型
					data.Entities[entity.Name] = entity.Value.(string)
				case int:
					//int类型
					data.Entities[entity.Name] = strconv.Itoa(entity.Value.(int))
				case float64:
					data.Entities[entity.Name] = strconv.Itoa(int(entity.Value.(float64)))
				default:
					log.Printf("数据 %v 既不是string也不是int类型, 而为 %T", entity.Value, entity.Value)
				}
				// data.Entities[entity.Name] = entity.Value
			}
		}
	}
	if len(msg.Intents) == 0 {
		// return fmt.Sprintf("呜呜, 我听不懂你在说什么. 但我猜:%s", Traits)
		data.Intents = 0
		return data, nil
	}
	// return fmt.Sprintf("我认为你在 %s, 并且可能具有以下特征:%s", msg.Intents[0].Name, Traits)
	intents := map[string]uint8{
		"Good_morning": 1,
		"Good_night":   2,
		"Humiliate":    3,
		"Praise":       4,
		"Stating":      5,
		"Awaken":       6,
		"Wish":         7,
		"Reminder":     8,
	}
	data.Intents = intents[msg.Intents[0].Name]
	return data, nil
}
