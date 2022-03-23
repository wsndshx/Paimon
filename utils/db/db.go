package db

import (
	"encoding/binary"
	"errors"
	"os"
	"path"
)

const (
	Del = iota
	Put
)

type KVdb struct {
	index    map[string]uint64 // 索引信息
	fileInfo File              // 数据库文件
}

type Entry struct {
	KeySize   uint32 // 记录键的长度
	ValueSize uint32 // 记录值的长度
	Flag      uint16 // 标记记录类型
	Key       []byte // 记录键数据
	Value     []byte // 记录值数据
}

type File struct {
	file *os.File // 文件的对象
	end  uint16   // 文件的尾巴
}

func (db *KVdb) Close() error {
	return db.fileInfo.file.Close()
}

// Crop 删除无效记录, 压缩数据库大小
func (db *KVdb) Crop() error {
	oldName := db.fileInfo.file.Name()
	newName := path.Dir(oldName) + "/tempDB"
	// 首先更新一次索引
	db.loadIndex()
	// 创建新的数据库文件
	newDB, _ := NewDB(newName)
	// 将索引中的内容写入新数据库
	for k, v := range db.index {
		value, err := db.getData(k, v)
		if err != nil {
			return errors.New("获取值时出现异常: " + err.Error())
		}
		newDB.Put([]byte(k), value)
	}
	newDB.Close()

	// 删除旧数据库
	db.Close()
	if err := os.Remove(oldName); err != nil {
		return errors.New("删除旧数据库文件时出现异常: " + err.Error())
	}
	// 将新数据库文件改名
	os.Rename(newName, oldName)
	return nil
}

// Get 获取数据库中指定键的值
func (db *KVdb) Get(key string) ([]byte, error) {
	// 检查索引
	index, isOK := db.index[key]
	if !isOK {
		return nil, errors.New("该值不存在")
	}

	content, err := db.getData(key, index)
	if err != nil {
		return nil, errors.New("获取值时出现异常: " + err.Error())
	}

	return content, nil
}

// getData 获取某条记录的值
func (db *KVdb) getData(key string, index uint64) ([]byte, error) {
	// 读取头部信息
	head := make([]byte, 10)
	if _, err := db.fileInfo.file.ReadAt(head, int64(index)); err != nil {
		return nil, err
	}
	// 解码头部
	keySize := binary.BigEndian.Uint32(head[0:4])
	valueSize := binary.BigEndian.Uint32(head[4:8])
	// Flag := binary.BigEndian.Uint16(head[8:10])

	// 读取内容
	content := make([]byte, keySize+valueSize)
	if _, err := db.fileInfo.file.ReadAt(content, int64(index+10)); err != nil {
		return nil, err
	}
	if string(content[0:keySize]) != key {
		return nil, errors.New("需要读取的键为\"" + key + "\", 然而读取到了\"" + string(content[0:keySize]) + "\"")
	}

	return content[keySize:], nil
}

// Put 向数据库中写入指定键值对
func (db *KVdb) Put(key []byte, value []byte) error {
	// 从结尾写入
	db.fileInfo.file.Seek(0, 2)
	index := db.fileInfo.end
	entryLen := 10 + len(key) + len(value)
	// 将需要写入的数据转为[]byte
	entry := make([]byte, int64(entryLen))
	binary.BigEndian.PutUint32(entry[0:4], uint32(len(key)))
	binary.BigEndian.PutUint32(entry[4:8], uint32(len(value)))
	binary.BigEndian.PutUint16(entry[8:10], Put)
	copy(entry[10:10+len(key)], key)
	copy(entry[10+len(key):], value)

	// 将数据写入文件
	if _, err := db.fileInfo.file.Write(entry); err != nil {
		return err
	}

	// 更新索引
	db.index[string(key)] = uint64(index)
	db.fileInfo.end += uint16(entryLen)

	return nil
}

// Open 打开指定数据库
func Open(fileName string) (*KVdb, error) {
	db := &KVdb{
		index: map[string]uint64{},
	}
	var err error

	// 判断是否存在该文件
	if info, err := os.Stat(fileName); err != nil {
		return nil, errors.New("指定数据库文件不存在: " + err.Error())
	} else {
		db.fileInfo.end = uint16(info.Size())
	}

	// 打开文件
	if db.fileInfo.file, err = os.OpenFile(fileName, os.O_RDWR, 0644); err != nil {
		return nil, errors.New("指定数据库文件打开失败: " + err.Error())
	}

	// 加载索引
	db.loadIndex()
	return db, nil
}

// loadIndex 加载索引
func (db *KVdb) loadIndex() error {
	var index int64 = 0
	// 从第一条开始读取至最后
	for {
		// 读取头部
		head := make([]byte, 10)
		if _, err := db.fileInfo.file.ReadAt(head, index); err != nil {
			return err
		}

		keySize := binary.BigEndian.Uint32(head[0:4])
		valueSize := binary.BigEndian.Uint32(head[4:8])
		Flag := binary.BigEndian.Uint16(head[8:10])

		// 读取内容
		data := make([]byte, keySize+valueSize)
		if _, err := db.fileInfo.file.ReadAt(data, index+10); err != nil {
			return err
		}
		if Flag == Del {
			// 在索引中删除
			delete(db.index, string(data[0:keySize]))
			continue
		}
		// 加入索引
		db.index[string(data[0:keySize])] = uint64(index)

		index = index + int64(keySize+valueSize+10)
		if index == int64(db.fileInfo.end) {
			break
		}
	}

	return nil
}

// NewDB 创建名为fileName的数据库, 并返回该数据库对象
func NewDB(fileName string) (*KVdb, error) {
	// 判断是否存在相同的文件/目录
	if _, err := os.Stat(fileName); err == nil {
		// 文件存在, 停止创建
		return nil, errors.New("文件名/目录存在, 无法创建数据库文件")
	}

	db := KVdb{
		index: make(map[string]uint64),
		fileInfo: File{
			end: 0,
		},
	}

	// 创建新的数据库文件
	{
		var err error
		if db.fileInfo.file, err = os.Create(fileName); err != nil {
			return nil, err
		}
	}

	return &db, nil
}
