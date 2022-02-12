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
}

// cqhttp的通信地址
var Host string

type Reply struct {
	Message_type string `json:"message_type"`
	User_id      int64  `json:"user_id"`
	Group_id     int64  `json:"group_id"`
	Message      string `json:"message"`
}

func init() {
	log.SetPrefix("[HTTP]")
	log.SetFlags(0)
}

// Reply 回复消息
func (message Reply) Reply() {
	type reply struct {
		Message_id int32 `json:"message_id"`
	}
	Json := reply{}
	JsonBody := new(bytes.Buffer)
	json.NewEncoder(JsonBody).Encode(message)

	res, err := http.Post(Host+"/send_msg", "application/json;charset=utf-8", JsonBody)
	if err != nil {
		log.Println(err.Error())
	}
	jsonData, _ := io.ReadAll(res.Body)
	json.Unmarshal(jsonData, &Json)
	log.Printf("发送消息`%s`成功, 返回消息ID: %d\n", message.Message, Json.Message_id)
	defer res.Body.Close()
}
