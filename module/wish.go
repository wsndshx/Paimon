package module

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wsndshx/Paimon/utils"
)

// 这里是模拟抽卡的噢....

// 先整个常规池
var ResidentRole []string = []string{
	"刻晴",
	"莫娜",
	"七七",
	"迪卢克",
	"琴",
	"阿莫斯之弓",
	"天空之翼",
	"四风原典",
	"天空之卷",
	"和璞鸢",
	"天空之脊",
	"狼的末路",
	"天空之傲",
	"天空之刃",
	"风鹰剑",
	"云堇",
	"九条裟罗",
	"五郎",
	"早柚",
	"托马",
	"烟绯",
	"罗莎莉亚",
	"辛焱",
	"砂糖",
	"迪奥娜",
	"重云",
	"诺艾尔",
	"班尼特",
	"菲谢尔",
	"凝光",
	"行秋",
	"北斗",
	"香菱",
	"安柏",
	"雷泽",
	"凯亚",
	"芭芭拉",
	"丽莎",
	"弓藏",
	"祭礼弓",
	"绝弦",
	"西风猎弓",
	"昭心",
	"祭礼残章",
	"流浪乐章",
	"西风秘典",
	"西风长枪",
	"匣里灭辰",
	"雨裁",
	"祭礼大剑",
	"钟剑",
	"西风大剑",
	"匣里龙吟",
	"祭礼剑",
	"笛剑",
	"西风剑",
	"弹弓",
	"神射手之誓",
	"鸦羽弓",
	"翡玉法球",
	"讨龙英杰谭",
	"魔导绪论",
	"黑缨枪",
	"以理服人",
	"沐浴龙血的剑",
	"铁影阔剑",
	"飞天御剑",
	"黎明神剑",
	"冷刃",
}

var times_Golden uint8 = 0
var times_Purple uint8 = 0

// wish 用于获取祈愿得到的物品星级
//
// 0 - 三星
// 1 - 四星
// 2 - 五星
func wish() uint8 {
	// 基础概率
	// 五星
	// 获取概率为0.6%; 在命中5星时, 50%概率为本期up, 否则下次五星必为本期up
	// 四星
	// 获取概率为5.1%, 角色占2.55%, 武器占2.55%; 在命中四星时, 50%概率为本期up, 否则下次四星必为本期up
	// 保底机制
	// 五星
	// 在90抽内必须出一次五星, 出后重置;
	// 四星
	// 每10次必出四星或以上物品, 四星的概率为99.4%, 五星的概率为0.6%

	// 设置默认的出金概率
	fiveProbability := 6
	// 设置默认的出紫概率
	fourProbability := 6 + 51

	// 抽取的次数+1
	times_Golden++
	times_Purple++
	// 创建随机数(0-1000)
	boom := rand.Intn(1000)

	// 若连续90次未出金, 则触发保底
	// if times_Golden == 90 {
	// 	fiveProbability = 1000
	// 	fourProbability = 1000
	// 	times_Golden = 0
	// }
	if times_Golden > 73 {
		switch times_Golden {
		case 74:
			fiveProbability = 66
		case 75:
			fiveProbability = 126
		case 76:
			fiveProbability = 186
		case 77:
			fiveProbability = 246
		case 78:
			fiveProbability = 306
		case 79:
			fiveProbability = 366
		case 80:
			fiveProbability = 426
		case 81:
			fiveProbability = 486
		case 82:
			fiveProbability = 546
		case 83:
			fiveProbability = 606
		case 84:
			fiveProbability = 666
		case 85:
			fiveProbability = 726
		case 86:
			fiveProbability = 786
		case 87:
			fiveProbability = 846
		case 88:
			fiveProbability = 906
		case 89:
			fiveProbability = 966
		case 90:
			fiveProbability = 1000
		}
		fourProbability += fiveProbability
	}

	// 若连续10次未出紫, 则提升概率
	// if times_Purple >= 10 {
	// 	fourProbability = 1000 - fiveProbability
	// 	times_Purple = 0
	// }
	if fiveProbability > 439 {
		if times_Purple >= 9 {
			fourProbability = 1000
		}
	} else if fiveProbability > 6 {
		fourProbability += fiveProbability
	} else {
		if times_Purple == 9 {
			fourProbability = 561
		} else if times_Purple > 9 {
			fourProbability = 1000
		}
	}

	// 开始抽取
	if boom > fourProbability-1 {
		// 这里是三星
		return 0
	} else if boom > fiveProbability-1 {
		// 这里是四星
		times_Purple = 0
		return 1
	} else {
		// 这里是五星
		times_Golden = 0
		return 2
	}
}

// Resident 常规祈愿
func Resident(times int) (result []string, err error) {
	// 开始祈愿噢~
	// 设置一个随机数
	rand.Seed(time.Now().Unix())
	{
		// 这里判断一下需要执行的次数, 当执行次数大于20时关闭文本输出, 转而使用网页显示
		if times <= 20 {
			for i := 0; i < times; i++ {
				switch wish() {
				case 0:
					result = append(result, "(蓝)"+ResidentRole[rand.Intn(13)+55])
				case 1:
					result = append(result, "(紫)"+ResidentRole[rand.Intn(40)+15])
				case 2:
					result = append(result, "(金)"+ResidentRole[rand.Intn(15)])
				}
			}
		} else {
			// 获取抽取结果并暂存
			type data struct {
				Name  string
				Grade string
				Type  string
			}
			var datas struct {
				Data   []data
				Golden uint
				Purple uint
			}
			datas.Golden = 0
			datas.Purple = 0
			for i := 0; i < times; i++ {
				switch wish() {
				case 0:
					datas.Data = append(datas.Data, data{
						Name:  ResidentRole[rand.Intn(13)+56],
						Grade: "三星",
						Type:  "武器",
					})
				case 1:
					datas.Purple++
					Type := ""
					Name := ""
					switch rand.Intn(2) {
					case 0:
						Type = "角色"
						Name = ResidentRole[rand.Intn(23)+15]
					case 1:
						Type = "武器"
						Name = ResidentRole[rand.Intn(18)+38]
					}
					datas.Data = append(datas.Data, data{
						Name:  Name,
						Grade: "四星",
						Type:  Type,
					})
				case 2:
					datas.Golden++
					Type := ""
					Name := ""
					switch rand.Intn(2) {
					case 0:
						Type = "角色"
						Name = ResidentRole[rand.Intn(5)]
					case 1:
						Type = "武器"
						Name = ResidentRole[rand.Intn(10)+5]
					}
					datas.Data = append(datas.Data, data{
						Name:  Name,
						Grade: "五星",
						Type:  Type,
					})
				}
			}
			// id 页面id或者数据库id
			var id string

			// 创建结果页面
			{
				pageJson := utils.NewPage{}
				pageJson.Parent.Type = "page_id"
				pageJson.Parent.Page_id = "3fae29a9c883424585bc2aebc3508487"
				pageJson.Properties = json.RawMessage(
					`{"title":{"title":[{"text":{"content":"祈愿结果 ` + time.Now().Format("2006-01-02 15:04:05") + `"}}]}}`)
				pageJson.Children = json.RawMessage(
					fmt.Sprintf(`[{"type":"paragraph","paragraph":{"rich_text":[{"type":"text","text":{"content":"本次抽取的次数: %d"}}]}},{"type":"paragraph","paragraph":{"rich_text":[{"type":"text","text":{"content":"其中包含"}},{"type":"text","text":{"content":" %d "},"annotations":{"color":"yellow_background"}},{"type":"text","text":{"content":"个五星物品，"}},{"type":"text","text":{"content":" %d "},"annotations":{"color":"purple_background"}},{"type":"text","text":{"content":"个四星物品"}}]}}]`, times, datas.Golden, datas.Purple))
				if id, err = utils.PostPage(pageJson); err != nil {
					return nil, fmt.Errorf("创建新页面失败: %s", err)
				} else if id == "" {
					return nil, fmt.Errorf("创建新页面失败: 未知错误, 返回页面ID为空")
				}
			}
			// 存入结果页面
			result = append(result, "https://paimonnya.notion.site/"+id)
			// 创建数据库
			if id, err = utils.PostDatabase(id); err != nil {
				return nil, fmt.Errorf("创建数据库失败: %s", err)
			} else if id == "" {
				return nil, fmt.Errorf("创建数据库失败: 未知错误, 返回数据库ID为空")
			}
			// 上传数据到结果页面
			go func() {
				databasePage := utils.NewDataPage{}
				databasePage.Parent.Type = "database_id"
				databasePage.Parent.Database_id = id
				for _, d := range datas.Data {
					err_times := 0
					// 构建数据库页面
					databasePage.Properties = json.RawMessage(fmt.Sprintf(`{"类型":{"type":"select","select":{"name":"%s"}},"等级":{"type":"select","select":{"name":"%s"}},"名称":{"type":"title","title":[{"type":"text","text":{"content":"%s"}}]}}`, d.Type, d.Grade, d.Name))
				reset:
					// 存入数据
					if _, err = utils.PostPage(databasePage); err != nil {

						if err_times >= 3 {
							log.Printf("存入数据库时出错: %s, 连续错误次数达3次, 退出当前任务...", err)
							return
						}
						// 这里加一个报错吧
						err_times++
						log.Printf("存入数据库时出错: %s, 将在%ds后重试", err, err_times*3)
						time.Sleep(time.Duration(3*err_times) * time.Second)
						goto reset
					}
					// 阻塞250毫秒, 避免触发429
					time.Sleep(250 * time.Millisecond)
				}
			}()
			return result, nil
			// 上传数据到全局数据库页面(待完成)
		}
	}
	return result, nil
}
