package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func init() {
	UserList = make(map[uint64]string)
}

type NewPage struct {
	Parent struct {
		Page_id string `json:"page_id"`
	} `json:"parent"`
	Properties interface{} `json:"properties"`
	Children   interface{} `json:"children"`
}

type NewDataPage struct {
	Parent struct {
		Database_id string `json:"database_id"`
	} `json:"parent"`
	Properties interface{} `json:"properties"`
}

var (
	Notion_token     string
	Wish_database_id string
	Wish_result_id   string
	UserList         map[uint64]string
)

// 获取指定QQ号对应的名称(如果有)
func GetUser(id uint64) string {
	if name, ok := UserList[id]; ok {
		return name
	}
	return fmt.Sprintf("%d", id)
}

// 向指定块中添加子块
func AddChildren(block_id string) []byte {
	url := "https://api.notion.com/v1/blocks/" + block_id + "/children"
	newChildren := struct {
		Children interface{} `json:"children"`
	}{}
	newChildren.Children = json.RawMessage(`[{"object":"block","type":"child_database","child_database":{"title":"这是一个标题","properties":{"等级":{"name":"等级","type":"select","select":{"options":[{"name":"五星","color":"yellow"},{"name":"四星","color":"purple"},{"name":"三星","color":"blue"}]}},"类型":{"name":"类型","type":"select","select":{"options":[{"name":"武器","color":"red"},{"name":"角色","color":"green"}]}},"名称":{"id":"title","name":"名称","type":"title","title":{}}}}}]`)
	newChildrenJson, _ := json.Marshal(newChildren)
	req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(newChildrenJson)))

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", "2022-02-22")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Notion_token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

// 获取指定块的内容
func GetBlockChildren(block_id string) []byte {
	url := "https://api.notion.com/v1/blocks/" + block_id + "/children?page_size=100"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", "2022-02-22")
	req.Header.Add("Authorization", "Bearer "+Notion_token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

// 获取指定块的信息
func GetBlock(block_id string) []byte {
	url := "https://api.notion.com/v1/blocks/" + block_id

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", "2022-02-22")
	req.Header.Add("Authorization", "Bearer "+Notion_token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

// 获取指定页面的信息
func GetPage(page_id string) []byte {
	url := "https://api.notion.com/v1/pages/" + page_id

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", "2022-02-22")
	req.Header.Add("Authorization", "Bearer "+Notion_token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func (data NewPage) PostPage() (string, error) {
	return postPage(data)
}
func (data NewDataPage) PostPage() (string, error) {
	return postPage(data)
}

// 创建页面, 返回创建的新页面的ID
func postPage(data interface{}) (string, error) {
	var res *http.Response
	{
		var err error
		var newPageJson []byte
		var payload *strings.Reader

		if newPageJson, err = json.Marshal(data); err != nil {
			return "", err
		}
		payload = strings.NewReader(string(newPageJson))

		request := NotionRequest{
			API:  "/v1/pages",
			Type: "POST",
			Body: payload,
			Res:  make(chan *http.Response),
		}
		requestQueue.Queue <- request
		res = <-request.Res
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)

	// 只提取出新创建的页面的id
	bodyJson := struct {
		Message string `json:"message"`
		Id      string `json:"id"`
	}{}
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("%s: %s", res.Status, bodyJson.Message)
	}
	bodyJson.Id = strings.Replace(bodyJson.Id, "-", "", -1)

	return bodyJson.Id, nil
}

// 查询数据库
func GetDatabase(parent string) []byte {
	url := "https://api.notion.com/v1/databases/" + parent

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", "2022-02-22")
	req.Header.Add("Authorization", "Bearer "+Notion_token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

// PostDatabase 用于在指定父页面上创建数据库页面
func PostDatabase(parent string) (string, error) {
	var res *http.Response
	{
		newDatabase := struct {
			Parent struct {
				Type    string `json:"type"`
				Page_id string `json:"page_id"`
			} `json:"parent"`
			Title      interface{} `json:"title"`
			Properties interface{} `json:"properties"`
		}{}
		newDatabase.Parent.Type = "page_id"
		newDatabase.Parent.Page_id = parent
		newDatabase.Title = json.RawMessage(`[{"type":"text","text":{"content":"详细数据"}}]`)
		newDatabase.Properties = json.RawMessage(`{"等级":{"name":"等级","type":"select","select":{"options":[{"name":"五星","color":"yellow"},{"name":"四星","color":"purple"},{"name":"三星","color":"blue"}]}},"类型":{"name":"类型","type":"select","select":{"options":[{"name":"武器","color":"red"},{"name":"角色","color":"green"}]}},"名称":{"id":"title","name":"名称","type":"title","title":{}}}`)
		newDatabaseJson, _ := json.Marshal(newDatabase)

		payload := strings.NewReader(string(newDatabaseJson))

		request := NotionRequest{
			API:  "/v1/databases",
			Type: "POST",
			Body: payload,
			Res:  make(chan *http.Response),
		}
		requestQueue.Queue <- request
		res = <-request.Res
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	/* 完整的解析结构体
	type AutoGenerated struct {
		Object      string      `json:"object"`
		ID          string      `json:"id"`
		Cover       interface{} `json:"cover"`
		Icon        interface{} `json:"icon"`
		CreatedTime time.Time   `json:"created_time"`
		CreatedBy   struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"created_by"`
		LastEditedBy struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"last_edited_by"`
		LastEditedTime time.Time `json:"last_edited_time"`
		Title          []struct {
			Type string `json:"type"`
			Text struct {
				Content string      `json:"content"`
				Link    interface{} `json:"link"`
			} `json:"text"`
			Annotations struct {
				Bold          bool   `json:"bold"`
				Italic        bool   `json:"italic"`
				Strikethrough bool   `json:"strikethrough"`
				Underline     bool   `json:"underline"`
				Code          bool   `json:"code"`
				Color         string `json:"color"`
			} `json:"annotations"`
			PlainText string      `json:"plain_text"`
			Href      interface{} `json:"href"`
		} `json:"title"`
		Properties struct {
			NAMING_FAILED struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Type   string `json:"type"`
				Select struct {
					Options []struct {
						ID    string `json:"id"`
						Name  string `json:"name"`
						Color string `json:"color"`
					} `json:"options"`
				} `json:"select"`
			} `json:"等级"`
			NAMING_FAILED struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Type   string `json:"type"`
				Select struct {
					Options []struct {
						ID    string `json:"id"`
						Name  string `json:"name"`
						Color string `json:"color"`
					} `json:"options"`
				} `json:"select"`
			} `json:"类型"`
			NAMING_FAILED struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Type  string `json:"type"`
				Title struct {
				} `json:"title"`
			} `json:"名称"`
		} `json:"properties"`
		Parent struct {
			Type   string `json:"type"`
			PageID string `json:"page_id"`
		} `json:"parent"`
		URL      string `json:"url"`
		Archived bool   `json:"archived"`
	}
	*/

	// 解析创建的数据库id
	bodyJson := struct {
		Id string `json:"id"`
	}{}
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		return "", fmt.Errorf("解析Json错误: %s", err)
	}
	bodyJson.Id = strings.Replace(bodyJson.Id, "-", "", -1)

	return bodyJson.Id, nil
}
