package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/robfig/cron/v3"
	"github.com/wsndshx/Paimon/utils/db"
)

var (
	cronDB *db.KVdb
)

type Cron struct {
	timer    *cron.Cron
	cronList []CronTask
}

type CronTask struct {
	Time     string
	Content  string
	TargetId int64
	EntryID  int
}

// NewCron 创建一个定时器
func NewCron() *Cron {
	c := &Cron{
		timer:    cron.New(),
		cronList: []CronTask{},
	}
	// 启动定时器
	c.timer.Start()
	return c
}

// CronLocal 加载数据库中的定时器记录
func (Cron *Cron) CronLocal() (err error) {
	// 从数据库中加载过去的任务
	// 判断是否存在旧的数据库文件
	if _, err = os.Stat("Cron.db"); err != nil {
		if os.IsNotExist(err) {
			// 文件不存在, 创建新的数据库文件
			cronDB, err = db.NewDB("Cron.db")
			return
		}
		// 文件打开失败
		return
	}
	// 读取数据库中的内容
	cronDB, err = db.Open("Cron.db")
	data, err := cronDB.Get("CronList")
	if err != nil {
		return
	}
	var history []CronTask
	err = decode(data, &history)

	// 重新添加任务
	for _, ct := range history {
		err = Cron.CronAdd(ct)
		if err != nil {
			return
		}
	}
	return
}

// CronList 列出当前存在的定时器数据
func (Cron *Cron) CronList() string {
	List := "------"
	for _, ct := range Cron.cronList {
		List += fmt.Sprintf("任务ID: %d\n表达式: %s\n内容: %s\n目标: %d\n------", ct.EntryID, ct.Time, ct.Content, ct.TargetId)
	}
	return fmt.Sprint(List)
}

// CronAdd 添加定时器
func (Cron *Cron) CronAdd(in CronTask) error {
	// 先注册
	ID, err := Cron.timer.AddFunc(in.Time, func() {
		log.Println("执行定时任务......")
		msg := Reply{
			Message_type: "group",
			Group_id:     in.TargetId,
			Message:      in.Content,
		}
		msg.Reply()

		// 这里删除数据库的内容
	})
	if err == nil {
		in.EntryID = int(ID)
		// 将内容写入列表
		Cron.cronList = append(Cron.cronList, in)
	}

	return err
}

// CronClose 在程序结束时调用, 把未完成的任务存入数据库
func (Cron *Cron) CronClose() error {
	data, err := encode(Cron.cronList)
	if err != nil {
		return err
	}
	cronDB.Put([]byte("CronList"), data)
	return nil
}

// encode 用于编码写入数据库的内容
func encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 用于解码从数据库中得到的内容
func decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
