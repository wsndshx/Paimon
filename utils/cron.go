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
	cronDB  *db.KVdb
	entryID = 0                          // 任务id计数器
	idToId  = make(map[int]cron.EntryID) //将entryID映射至cronID
)

type Cron struct {
	dbPATH   string
	timer    *cron.Cron
	cronList cronList
}

type cronList struct {
	index map[cron.EntryID]int
	list  []CronTask
}

type CronTask struct {
	Time     string
	Content  string
	TargetId int64
	EntryID  cron.EntryID
}

// NewCron 创建一个定时器;
//
// 该函数接受一个可选的参数dbPath，用于指定数据库的存放路径
func NewCron(dbPath ...string) *Cron {
	// 缺省路径
	path := "data/Cron.db"
	if dbPath != nil {
		path = dbPath[0]
	}
	c := &Cron{
		dbPATH: path,
		timer:  cron.New(),
		cronList: cronList{
			index: make(map[cron.EntryID]int),
			list:  []CronTask{},
		},
	}
	c.timer.Start()
	return c
}

func (Cron *Cron) Remove(ID int) {
	// 删除定时器
	Cron.timer.Remove(cron.EntryID(ID))
	// 删除列表
	Cron.cronList.remove(cron.EntryID(ID))
	// 更新索引
	Cron.cronList.updateIndex()
}

func (cronList *cronList) remove(ID cron.EntryID) {
	// 将该任务从列表中删除
	log.Println("任务完成, 开始将该任务从列表中删除")
	if len(cronList.list) == 1 {
		cronList.list = []CronTask{}
		cronList.index = make(map[cron.EntryID]int)
		log.Println("删除成功")
		return
	}
	// 查找索引
	index := cronList.index[ID]
	// 若索引失败
	if cronList.list[index].EntryID != ID {
		log.Println("索引失败, 开始遍历")
		for i := range cronList.list {
			if cronList.list[i].EntryID == ID {
				cronList.list = append(cronList.list[:i], cronList.list[i+1:]...)
				return
			}
		}
		// 这里是没找到的情况
		log.Printf("删除任务失败: 未找到相应的任务ID: %d", ID)
		return
	}
	// 成功则直接移除当前元素
	cronList.list = append(cronList.list[:index], cronList.list[index+1:]...)
	log.Println("删除成功")
}

// Local 加载数据库中的定时器记录
func (Cron *Cron) Local() (err error) {
	// 从数据库中加载过去的任务
	// 判断是否存在旧的数据库文件
	if _, err = os.Stat(Cron.dbPATH); err != nil {
		if os.IsNotExist(err) {
			// 文件不存在, 创建新的数据库文件
			cronDB, err = db.NewDB(Cron.dbPATH)
			return
		}
		// 文件打开失败
		return
	}
	// 读取数据库中的内容
	cronDB, err = db.Open(Cron.dbPATH)
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

// List 列出当前存在的定时器数据
func (Cron *Cron) List() string {
	List := "------"
	for _, ct := range Cron.cronList.list {
		List += fmt.Sprintf("\n任务ID: %d\n表达式: %s\n内容: %s\n目标: %d\n------", ct.EntryID, ct.Time, ct.Content, ct.TargetId)
	}
	return fmt.Sprint(List)
}

// CronAdd 添加定时器
func (Cron *Cron) CronAdd(in CronTask) error {
	log.Println("开始添加定时任务......")
	Group_id := in.TargetId
	Message := in.Content
	// 先注册
	ID, err := Cron.timer.AddFunc(in.Time, func() {
		log.Println("执行定时任务......")
		msg := Reply{
			Message_type: "group",
			Group_id:     Group_id,
			Message:      Message,
		}
		msg.Reply()

		// 将该任务从列表中删除
		Cron.cronList.remove(idToId[entryID])
		// 更新索引
		Cron.cronList.updateIndex()
	})
	if err != nil {
		return err
	}
	idToId[entryID] = ID
	entryID++

	// 将内容写入列表
	in.EntryID = ID
	Cron.cronList.list = append(Cron.cronList.list, in)
	Cron.cronList.index[ID] = len(Cron.cronList.list) - 1

	return nil
}

func (cronList *cronList) updateIndex() {
	// 更新索引
	if len(cronList.list) == 0 {
		cronList.index = make(map[cron.EntryID]int)
		return
	}
	for i := range cronList.list {
		cronList.index[cronList.list[i].EntryID] = i
	}
}

// Close 在程序结束时调用, 把未完成的任务存入数据库
func (Cron *Cron) Close() error {
	log.Println("开始存储定时器任务")
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
