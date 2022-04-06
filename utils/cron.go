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
	cronList struct {
		index map[int]int
		list  []CronTask
	}
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
		timer: cron.New(),
		cronList: struct {
			index map[int]int
			list  []CronTask
		}{
			index: make(map[int]int),
			list:  []CronTask{},
		},
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
	for _, ct := range Cron.cronList.list {
		List += fmt.Sprintf("\n任务ID: %d\n表达式: %s\n内容: %s\n目标: %d\n------", ct.EntryID, ct.Time, ct.Content, ct.TargetId)
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

		// 将该任务从列表中删除
		log.Println("任务完成, 开始将该任务从列表中删除")
		if len(Cron.cronList.list) == 1 {
			Cron.cronList.list = []CronTask{}
			Cron.cronList.index = make(map[int]int)
			log.Println("删除成功")
			Cron.updateIndex()
			return
		}
		index := Cron.cronList.index[in.EntryID]
		if Cron.cronList.list[index].EntryID != in.EntryID {
			log.Println("索引失败, 开始遍历")
			for i, ct := range Cron.cronList.list {
				if ct.Time == in.Time && ct.Content == in.Content && ct.TargetId == in.TargetId {
					Cron.cronList.list = append(Cron.cronList.list[:i], Cron.cronList.list[i+1:]...)
					Cron.updateIndex()
					return
				}
			}
			// 这里是没找到的情况
			log.Printf("任务列表删除已完成任务时失败: 未找到相应的任务: \n任务ID(猜测): %d\n表达式: %s\n内容: %s\n目标: %d\n", in.EntryID, in.Time, in.Content, in.TargetId)
		}
		// 若索引成功则直接移除当前元素
		Cron.cronList.list = append(Cron.cronList.list[:index], Cron.cronList.list[index+1:]...)
		log.Println("删除成功")
		Cron.updateIndex()
	})
	if err != nil {
		return err
	}
	in.EntryID = int(ID)
	// 将内容写入列表
	Cron.cronList.list = append(Cron.cronList.list, in)
	Cron.cronList.index[int(ID)] = len(Cron.cronList.list) - 1

	return nil
}

func (Cron *Cron) updateIndex() {
	// 更新索引
	if len(Cron.cronList.list) == 0 {
		Cron.cronList.index = make(map[int]int)
		return
	}
	for i := range Cron.cronList.list {
		Cron.cronList.index[Cron.cronList.list[i].EntryID] = i
	}
}

// CronClose 在程序结束时调用, 把未完成的任务存入数据库
func (Cron *Cron) CronClose() error {
	data, err := encode(Cron.cronList.list)
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
