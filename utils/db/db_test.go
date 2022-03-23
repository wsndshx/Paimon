package db_test

import (
	"testing"

	"github.com/wsndshx/Paimon/utils/db"
)

func TestNewDB(t *testing.T) {
	// 新建数据库
	testDB, err := db.NewDB("testDB")
	defer testDB.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestPut(t *testing.T) {
	// 打开数据库
	testDB, err := db.Open("testDB")
	defer testDB.Close()
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 10; i++ {
		// 写入数据
		if err := testDB.Put([]byte("KeyStr"), []byte("Value")); err != nil {
			t.Error(err)
		}
		if err := testDB.Put([]byte("NekoStr"), []byte("Miao")); err != nil {
			t.Error(err)
		}
		if err := testDB.Put([]byte("NekoStr"), []byte("MiaoMiao")); err != nil {
			t.Error(err)
		}
	}
}

func TestCrop(t *testing.T) {
	// 打开数据库
	testDB, err := db.Open("testDB")
	defer testDB.Close()
	if err != nil {
		t.Error(err)
		return
	}
	// 裁剪数据库
	if err := testDB.Crop(); err != nil {
		t.Error(err)
		return
	}
}

func TestGet(t *testing.T) {
	// 打开数据库
	testDB, err := db.Open("testDB")
	defer testDB.Close()
	if err != nil {
		t.Error(err)
		return
	}
	// 读取数据
	data, err := testDB.Get("KeyStr")
	if err != nil {
		t.Error(err)
		return
	}
	if string(data[:]) != "Value" {
		t.Errorf("数据读取异常, 预期返回值为 `Value`, 却得到了 `%s`", string(data[:]))
	}
	data, err = testDB.Get("NekoStr")
	if err != nil {
		t.Error(err)
		return
	}
	if string(data[:]) != "MiaoMiao" {
		t.Errorf("数据读取异常, 预期返回值为 `MiaoMiao`, 却得到了 `%s`", string(data[:]))
	}
}
