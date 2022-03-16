package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var requestQueue *RequestQueue

func init() {
	// 初始化任务队列
	requestQueue = &RequestQueue{
		Queue: make(chan NotionRequest, 256),
	}
	go requestQueue.notion()
}

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

type ReplyP struct {
	User_id int64  `json:"user_id"`
	Message string `json:"message"`
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
func (message ReplyP) Reply() []byte {
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
	body, _ := ioutil.ReadAll(res.Body)

	return body
}
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

// NotionRequest 用于构建请求
type NotionRequest struct {
	API  string              //需要请求的接口
	Type string              //请求的类型
	Body io.Reader           //发送出去的数据
	Res  chan *http.Response //返回的数据
}
type RequestQueue struct {
	Queue     chan NotionRequest // 普通队列
	HighQueue chan NotionRequest // 会被优先处理的队列
}

// notion 这里用来处理notion相关api的请求
func (req *RequestQueue) notion() {
	// 接受请求
	for {
		select {
		case data := <-req.HighQueue:
			processing(data)
		default:
			select {
			case data := <-req.HighQueue:
				processing(data)
			case data := <-req.Queue:
				processing(data)
			}
		}
	}
}

// processing 处理队列
func processing(data NotionRequest) {
	url := "https://api.notion.com" + data.API
	req, _ := http.NewRequest(data.Type, url, data.Body)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", "2022-02-22")
	req.Header.Add("Authorization", "Bearer "+Notion_token)
	if data.Body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	var res *http.Response
	var err error
	relieve := 0
retry:
	if res, err = http.DefaultClient.Do(req); err != nil {
		log.Println(err)
		return
	} else if res.StatusCode == 429 {
		relieve++
		log.Printf("notion api 请求过快! 将在%ds后重新启动队列\n", relieve)
		time.Sleep(time.Duration(relieve*3) * time.Second)
		goto retry
	}
	// 将数据返回
	data.Res <- res
	close(data.Res)
}
