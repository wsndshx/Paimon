package utils

import (
	"bytes"
	"encoding/json"
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

// Reply 回复消息
func (message Reply) Reply() {
	JsonBody := new(bytes.Buffer)
	json.NewEncoder(JsonBody).Encode(message)

	req, _ := http.NewRequest("POST", Host+"/send_msg", JsonBody)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	resp.Body.Close()
}
