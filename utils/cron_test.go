package utils_test

import (
	"testing"

	"github.com/wsndshx/Paimon/utils"
)

var testCron *utils.Cron

func init() {
	testCron = utils.NewCron("TestCron.db")
}

func TestCronLocal(t *testing.T) {
	if err := testCron.CronLocal(); err != nil {
		t.Error(err)
	}
}

func TestCronList(t *testing.T) {
	t.Log("目前包含的任务为: \n" + testCron.CronList())
}

func TestCronAdd(t *testing.T) {
	cron := utils.CronTask{
		Time:     "5 0 * 8 *",
		Content:  "这是一个提醒任务",
		TargetId: 0,
	}
	if err := testCron.CronAdd(cron); err != nil {
		t.Error(err)
	}
	t.Log("目前包含的任务为: \n" + testCron.CronList())
}

func TestCronClose(t *testing.T) {
	if err := testCron.CronClose(); err != nil {
		t.Error(err)
	}
}
