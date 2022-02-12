package group

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
)

const Host = "127.0.0.1:5700"

type Reply struct {
	Group_id    int64  `json:"group_id"`
	Message     string `json:"message"`
	Auto_escape bool   `json:"auto_escape"`
}

func Handle(message string, num int64) {
	// 匹配早安消息
	morning := regexp.MustCompile(`^(早安|早上好|早鸭) | (早安|早上好|早鸭)$`)
	// 匹配晚安消息
	night := regexp.MustCompile(`^(晚安|我睡了|睡) | (晚安|我睡了)$`)
	if morning.MatchString(message) {
		greeting(1, num)
	} else if night.MatchString(message) {
		greeting(2, num)
	}
}

func reply(num int64, message string) {
	Json := Reply{
		Group_id:    num,
		Message:     message,
		Auto_escape: false,
	}

	JsonBody := new(bytes.Buffer)
	json.NewEncoder(JsonBody).Encode(Json)

	req, _ := http.NewRequest("POST", Host, JsonBody)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	resp.Body.Close()
}
