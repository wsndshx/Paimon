package module

import (
	"math/rand"
	"time"
)

// 这里是模拟抽卡的噢....

// 先整个常规池
var ResidentRole []string = []string{
	"刻晴",
	"莫娜",
	"七七",
	"迪卢克",
	"琴",
	"天空之翼",
	"天空之卷",
	"天空之脊",
	"天空之傲",
	"天空之刃",
	"阿莫斯之弓",
	"四风原典",
	"和璞鸢",
	"狼的末路",
	"风鹰剑",
	"丽莎",
	"芭芭拉",
	"凯亚",
	"雷泽",
	"安柏",
	"香菱",
	"北斗",
	"行秋",
	"凝光",
	"菲谢尔",
	"班尼特",
	"诺艾尔",
	"重云",
	"砂糖",
	"迪奥娜",
	"辛焱",
	"罗莎莉亚",
	"烟绯",
	"早柚",
	"九条裟罗",
	"托马",
	"五郎",
	"西风猎弓",
	"西风秘典",
	"西风长枪",
	"西风大剑",
	"西风剑",
	"笛剑",
	"祭礼剑",
	"匣里龙吟",
	"钟剑",
	"祭礼大剑",
	"雨裁",
	"匣里灭辰",
	"流浪乐章",
	"祭礼残章",
	"昭心",
	"绝弦",
	"祭礼弓",
	"弓藏",
	"冷刃",
	"黎明神剑",
	"飞天御剑",
	"铁影阔剑",
	"沐浴龙血的剑",
	"以理服人",
	"黑缨枪",
	"魔导绪论",
	"讨龙英杰",
	"翡玉法球",
	"鸦羽弓神",
	"神射手之誓",
	"弹弓",
}

var times_Golden uint8 = 0
var times_Purple uint8 = 0

// wish 用于获取祈愿得到的物品星级
//
// 0 - 三星
// 1 - 四星
// 2 - 五星
func Wish() uint8 {
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
func Resident(times int) (result []string) {
	// 开始祈愿噢~
	// 设置一个随机数
	rand.Seed(time.Now().Unix())
	{
		for i := 0; i < times; i++ {
			switch Wish() {
			case 0:
				result = append(result, "(蓝)"+ResidentRole[rand.Intn(13)+55])
			case 1:
				result = append(result, "(紫)"+ResidentRole[rand.Intn(40)+15])
			case 2:
				result = append(result, "(金)"+ResidentRole[rand.Intn(15)])
			}
		}
	}
	return result
}
