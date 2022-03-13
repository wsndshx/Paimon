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

// 抽取结果
type datas struct {
	Data   []data // 每次抽取到的物品数据
	Golden uint8  // 抽取到的金色物品数量
	Purple uint8  // 抽取到的紫色物品数量
}

// 抽取到的物品数据
type data struct {
	Name  string // 物品名称
	Grade string // 物品等级
	Type  string // 物品类型
}

var (
	IsRoleUpPurple bool  = false // 下次祈愿获取的4星物品必定为本期4星UP角色
	IsRoleUpGolden bool  = false // 下次祈愿获取的5星物品必定为本期5星UP角色
	times_Golden   uint8 = 0     // 统计未抽到的金色物品的次数
	times_Purple   uint8 = 0     // 统计未抽到的紫色物品的次数
)

// 整个角色活动祈愿池
var UpRole map[string][]string = map[string][]string{
	"雷神": {
		"雷电将军",
		"九条裟罗",
		"辛焱",
		"班尼特",
	},
	"心海": {
		"珊瑚宫心海",
		"九条裟罗",
		"辛焱",
		"班尼特",
	},
}

// 先整个常规池
var ResidentWeapon []string = []string{
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
var ResidentRole []string = []string{
	"刻晴",
	"莫娜",
	"七七",
	"迪卢克",
	"琴",
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
}

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

// Role 角色活动祈愿
//
// times: 需要抽取的次数; qq: 进行祈愿的用户; up: 需要抽取的up池
func Role(times uint8, qq uint64, up string) (result []string, err error) {
	var upPool []string // 用户选择的角色池
	{
		var ok bool
		if upPool, ok = UpRole[up]; !ok {
			return nil, fmt.Errorf("%s: 不存在该角色池!", up)
		}
	}

	// 设置一个随机数
	rand.Seed(int64(qq))

	datas := datas{}
	datas.Golden = 0
	datas.Purple = 0

	// 开始祈愿噢~
	{
		var i uint8
		for i = 0; i < times; i++ {
			switch wish() {
			case 0:
				// 出三星物品
				datas.Data = append(datas.Data, data{
					Name:  ResidentWeapon[rand.Intn(13)+28],
					Grade: "三星",
					Type:  "武器",
				})

			case 1:
				// 出四星物品
				datas.Purple++
				Type := ""
				Name := ""
				if IsRoleUpPurple {
					// 下次祈愿获取的4星物品必定为本期4星UP角色
					Type = "角色"
					Name = upPool[rand.Intn(3)+1]
					goto breakPurpleSwitch
				}
				// 判断是否UP
				switch rand.Intn(2) {
				case 0:
					// 有50.000%的概率为本期4星UP角色
					Type = "角色"
					Name = upPool[rand.Intn(3)+1]
					IsRoleUpPurple = false
				case 1:
					// 下次祈愿获取的4星物品必定为本期4星UP角色
					IsRoleUpPurple = true
					// 常驻角色武器五五开
					switch rand.Intn(2) {
					case 0:
						Type = "角色"
						Name = ResidentRole[rand.Intn(23)+5]
					case 1:
						Type = "武器"
						Name = ResidentWeapon[rand.Intn(18)+10]
					}
				}

			breakPurpleSwitch:
				datas.Data = append(datas.Data, data{
					Name:  Name,
					Grade: "四星",
					Type:  Type,
				})
				// 出五星物品
			case 2:
				datas.Golden++
				Type := "角色"
				Name := ""
				if IsRoleUpGolden {
					Name = upPool[0]
					// 这里跳过下方的switch块
					goto breakGoldenSwitch
				}

				switch rand.Intn(2) {
				case 0:
					IsRoleUpGolden = true
					// 50%的概率获得非本期5星UP角色
					Name = ResidentRole[rand.Intn(5)]
				case 1:
					// 50%的概率为本期5星UP角色
					Name = upPool[0]
					IsRoleUpGolden = false
				}

			breakGoldenSwitch:
				datas.Data = append(datas.Data, data{
					Name:  Name,
					Grade: "五星",
					Type:  Type,
				})
			}
		}
	}

	return datas.getResult(times, qq)
}

// Resident 常规祈愿
func Resident(times uint8, qq uint64) (result []string, err error) {
	// 设置一个随机数
	rand.Seed(int64(qq))

	datas := datas{}
	datas.Golden = 0
	datas.Purple = 0

	// 开始祈愿噢~
	{
		var i uint8
		for i = 0; i < times; i++ {
			switch wish() {
			case 0:
				datas.Data = append(datas.Data, data{
					Name:  ResidentWeapon[rand.Intn(13)+28],
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
					Name = ResidentRole[rand.Intn(23)+5]
				case 1:
					Type = "武器"
					Name = ResidentWeapon[rand.Intn(18)+10]
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
					Name = ResidentWeapon[rand.Intn(10)]
				}
				datas.Data = append(datas.Data, data{
					Name:  Name,
					Grade: "五星",
					Type:  Type,
				})
			}
		}
	}

	return datas.getResult(uint8(times), qq)
}

// getResult 获取格式化的抽取结果
func (datas datas) getResult(times uint8, qq uint64) (result []string, err error) {
	// 上传数据到全局数据库页面
	go func() {
		databasePage := utils.NewDataPage{}
		databasePage.Parent.Database_id = utils.Wish_database_id

		// TODO 遍历datas.Data中包含的抽取结果,
		// 并将其上传至之前创建的数据库页面
		for _, data := range datas.Data {
			// 构建数据库页面
			databasePage.Properties = json.RawMessage(fmt.Sprintf(`{"等级":{"type":"select","select":{"name":"%s"}},"名称":{"type":"rich_text","rich_text":[{"type":"text","text":{"content":"%s"}}]},"时间":{"type":"date","date":{"start":"%s"}},"类型":{"type":"select","select":{"name":"%s"}},"操作人":{"type":"title","title":[{"type":"text","text":{"content":"%s"}}]}}`, data.Grade, data.Name, time.Now().Format("2006-01-02T15:04:05.000-07:00"), data.Type, utils.GetUser(qq)))
			// 统计错误次数
			err_times := 0
		reset:
			if _, err := databasePage.PostPage(); err != nil {
				// 错误处理
				if err_times > 3 {
					log.Printf("存入汇总数据库时出错: %s, 连续错误次数达3次, 退出当前任务...", err)
					return
				}
				err_times++
				log.Printf("存入汇总数据库时出错: %s, 将在%ds后重试", err, err_times*3)
				time.Sleep(time.Duration(3*err_times) * time.Second)
				goto reset
			}
		}
	}()

	// 这里判断一下祈愿的次数, 当祈愿次数大于等于20时关闭文本输出, 转而使用网页显示
	if times <= 20 {
		for _, d := range datas.Data {
			result = append(result, fmt.Sprintf("(%s)%s", d.Grade, d.Name))
		}
	} else {
		result = append(result, fmt.Sprintf("总抽取次数: %d\n金色: %d\n紫色: %d", times, datas.Golden, datas.Purple))
		// id 页面id或者数据库id
		var id string
		// 创建结果页面
		{
			pageJson := utils.NewPage{}
			pageJson.Parent.Page_id = utils.Wish_result_id
			pageJson.Properties = json.RawMessage(
				`{"title":{"title":[{"text":{"content":"祈愿结果 ` + time.Now().Format("2006-01-02 15:04:05") + `"}}]}}`)
			pageJson.Children = json.RawMessage(
				fmt.Sprintf(`[{"type":"paragraph","paragraph":{"rich_text":[{"type":"text","text":{"content":"本次抽取的次数: %d"}}]}},{"type":"paragraph","paragraph":{"rich_text":[{"type":"text","text":{"content":"其中包含"}},{"type":"text","text":{"content":" %d "},"annotations":{"color":"yellow_background"}},{"type":"text","text":{"content":"个五星物品，"}},{"type":"text","text":{"content":" %d "},"annotations":{"color":"purple_background"}},{"type":"text","text":{"content":"个四星物品"}}]}}]`, times, datas.Golden, datas.Purple))
			if id, err = pageJson.PostPage(); err != nil {
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
			databasePage.Parent.Database_id = id

			// TODO 遍历datas.Data中包含的抽取结果,
			// 并将其上传至之前创建的数据库页面
			for _, d := range datas.Data {
				// 构建数据库页面
				databasePage.Properties = json.RawMessage(fmt.Sprintf(`{"类型":{"type":"select","select":{"name":"%s"}},"等级":{"type":"select","select":{"name":"%s"}},"名称":{"type":"title","title":[{"type":"text","text":{"content":"%s"}}]}}`, d.Type, d.Grade, d.Name))
				// 统计错误次数
				err_times := 0
			reset:
				// 存入数据
				if _, err = databasePage.PostPage(); err != nil {
					if err_times >= 3 {
						log.Printf("存入结果页面数据库时出错: %s, 连续错误次数达3次, 退出当前任务...", err)
						return
					}
					// 这里加一个报错吧
					err_times++
					log.Printf("存入结果页面数据库时出错: %s, 将在%ds后重试", err, err_times*3)
					time.Sleep(time.Duration(3*err_times) * time.Second)
					goto reset
				}
			}
		}()
	}

	return result, nil
}
