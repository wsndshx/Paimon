package module

import (
	"fmt"
	"log"
	"strconv"

	"github.com/wsndshx/Paimon/utils"
)

var Timer *utils.Cron

func Cmd(cmd []string, msg *utils.Reply, user uint64) {
	errStr := ""
	isERR := false
	errInt := 0
OutCMD:
	switch cmd[0] {
	case "Reminder":
		if len(cmd) != 7 {
			isERR = true
			errStr = fmt.Sprintf("该命令允许的参数数量为6个, 然而却得到了%d个参数", (len(cmd) - 1))
			break OutCMD
		}
		cronStr := fmt.Sprintf("%s %s %s %s ?", cmd[2], cmd[3], cmd[4], cmd[5])
		targetId, err := strconv.ParseInt(cmd[1], 10, 64)
		if err != nil {
			isERR = true
			errStr = fmt.Sprintf("解析目标ID错误: %v\n", err)
			break OutCMD
		}

		// 添加定时器任务
		if err := Timer.CronAdd(utils.CronTask{
			Time:     cronStr,
			Content:  cmd[6],
			TargetId: targetId,
		}); err != nil {
			isERR = true
			errStr = fmt.Sprintf("添加定时器任务失败 : %v\n", err)
			break OutCMD
		}
		msg.Message = "派蒙记住了~"

	case "Wish":
		msg.Message = "少女祈祷中......"
		msg.Reply()

		// 存储祈愿结果
		var result []string

		// 解析抽取的次数
		var times uint8
		if time, err := strconv.ParseUint(cmd[3], 10, 8); err != nil {
			isERR = true
			errStr = "解析抽取次数失败: " + err.Error()
			break OutCMD
		} else {
			times = uint8(time)
		}

		// 解析二级指令
		switch cmd[1] {
		case "Resident":
			// 常驻祈愿
			if data, err := Resident(times, user); err != nil {
				isERR = true
				errStr = "常规祈愿失败: " + err.Error()
				break OutCMD
			} else {
				result = data
			}
		case "Role":
			// 角色祈愿
			if data, err := Role(times, user, cmd[2]); err != nil {
				isERR = true
				errStr = "角色活动祈愿失败: " + err.Error()
				break OutCMD
			} else {
				result = data
			}
		case "Weapon":
			// 武器祈愿
		}
		errInt++

		// 当祈愿次数小于等于20时使用纯文本展示,
		// 否则使用网站链接展示
		if times <= 20 {
			msg.Message = fmt.Sprintf("太好了旅行者, 抽到了这些东西呢: \n%v", result)
		} else {
			msg.Message = fmt.Sprintf("太好了旅行者, 抽到了这些东西呢: \n%s\n详细数据: %s", result[0], result[1])
		}

	case "cron":
		switch cmd[1] {
		case "list":
			errInt++
			msg.Message = "目前包含的任务为: \n" + Timer.List()
		case "remove":
			errInt++
			if len(cmd) != 3 {
				errInt++
				errStr = fmt.Sprintf("当前指令接受2个参数，却得到了%d个参数", len(cmd)-1)
				break
			}
			// 删除任务
			id, err := strconv.Atoi(cmd[2])
			if err != nil {
				isERR = true
				errStr = "解析任务id时出错: " + err.Error()
				errInt++
			}
			Timer.Remove(id)
			msg.Message = "任务删除成功"
		default:
			isERR = true
		}
		errInt++
	}
	errInt++
	if isERR {
		if errStr == "" {
			errStr = "不存在指令" + cmd[errInt]
		}
		errStr = fmt.Sprintf("指令在位置%d上解析错误: %s", errInt, errStr)

		log.Printf("指令执行错误: %s\n", errStr)
		msg.Message = "呜呜呜出错了: " + errStr
	}
	msg.Reply()
}
