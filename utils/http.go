package utils

import (
	"bytes"
	"encoding/json"
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
	JsonBody := new(bytes.Buffer)
	json.NewEncoder(JsonBody).Encode(message)

	res, err := http.Post(Host+"/send_msg", "application/json;charset=utf-8", JsonBody)
	if err != nil {
		log.Println(err.Error())
	}
	defer res.Body.Close()
}
