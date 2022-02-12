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

	req, err := http.NewRequest("POST", Host+"/send_msg", JsonBody)
	if err != nil {
		log.Println(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	resp.Body.Close()
}
