package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type HTTP interface {
	Reply()
	ChineseWS()
}

// cqhttp的通信地址
var Host string

type Reply struct {
	Message_type string `json:"message_type"`
	User_id      int64  `json:"user_id"`
	Group_id     int64  `json:"group_id"`
	Message      string `json:"message"`
}

type ChineseWS struct {
	Content string `json:"content"`
}

// Chinese Word Segmentation
func (message ChineseWS) ChineseWS() string {
	type reply struct {
		Data struct {
			Slices []string `json:"slices"`
		} `json:"data"`
	}
	slicesData := reply{}
	// slicesData := Slices{}
	JsonBody := new(bytes.Buffer)
	json.NewEncoder(JsonBody).Encode(message)
	res, err := http.Post(Host+"/.get_word_slices", "application/json;charset=utf-8", JsonBody)
	if err != nil {
		log.Println(err.Error())
	}
	// 解析返回值
	jsonData, _ := io.ReadAll(res.Body)
	json.Unmarshal(jsonData, &slicesData)
	if len(slicesData.Data.Slices) == 0 {
		log.Panicln("获取中文分词失败! 请检查`/.get_word_slices`是否正常")
	}
	log.Printf("获取`%s`的分词成功, 分词结果为%v\n", message.Content, slicesData.Data.Slices)

	defer res.Body.Close()
	slices := ""
	for _, v := range slicesData.Data.Slices {
		slices += v + " "
	}
	slices = slices[:len(slices)-1]
	return slices
}

// Reply 回复消息
func (message Reply) Reply() {
	// type reply struct {
	// 	Message_id int32 `json:"message_id"`
	// }
	// Json := reply{}
	JsonBody := new(bytes.Buffer)
	json.NewEncoder(JsonBody).Encode(message)

	res, err := http.Post(Host+"/send_msg", "application/json;charset=utf-8", JsonBody)
	if err != nil {
		log.Println(err.Error())
	} else {
		// jsonData, _ := io.ReadAll(res.Body)
		// json.Unmarshal(jsonData, &Json)
		log.Printf("发送消息`%s`成功\n", message.Message)
	}

	defer res.Body.Close()
}
