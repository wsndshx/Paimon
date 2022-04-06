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
	if err := testCron.Local(); err != nil {
		t.Error(err)
	}
}

func TestCronList(t *testing.T) {
	t.Log("目前包含的任务为: \n" + testCron.List())
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
	t.Log("目前包含的任务为: \n" + testCron.List())

	cron.Content = "这是应当被删除的记录"
	if err := testCron.CronAdd(cron); err != nil {
		t.Error(err)
	}
	t.Log("目前包含的任务为: \n" + testCron.List())
}

func TestRemove(t *testing.T) {
	testCron.Remove(2)
}

func TestCronClose(t *testing.T) {
	if err := testCron.Close(); err != nil {
		t.Error(err)
	}
}
