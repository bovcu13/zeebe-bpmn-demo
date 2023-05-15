package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
)

var token = GetAccessToken()

func main() {

	//variables := []map[string]string{
	//	{
	//		"name":  "meal",
	//		"value": "\"Chicken\"",
	//	},
	//}
	//CompleteTask("2251799813685266", variables)

	//GetTaskVariables("2251799813693348")

	//GetAllCreatedTasks()

	//fmt.Println(GetAccessToken())

	AssignTask("2251799813685257", "demo")
}

// GetAccessToken
// @Description: 取得Cookie
func GetAccessToken() string {
	client := resty.New()
	body := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {"demo"},
		"client_secret": {"M4SrdxOcznaNm9ZU9OWCb9htxHIdWAzv"},
	}

	resp, err := client.R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(body.Encode()).
		Post("http://localhost:18080/auth/realms/camunda-platform/protocol/openid-connect/token")
	if err != nil {
		fmt.Println(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		fmt.Println(err)
	}

	token, ok := data["access_token"].(string)
	if !ok {
		fmt.Println("Unable to get access token")
	}

	return "Bearer " + token
}

// GetCookie
// @Description: 取得Cookie
func GetCookie() (name, value string) {
	// post login form
	client := resty.New()

	resp, err := client.R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"username": "demo",
			"password": "demo",
		}).
		Post("http://localhost:8082/api/login")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.String())

	return resp.Cookies()[0].Name, resp.Cookies()[0].Value
}

// GetAllTasks
// @Description: 獲得所有任務詳細資訊
func GetAllTasks() []byte {
	client := resty.New()

	resp, err := client.R().SetHeader("Content-Type", "application/json").
		Post("http://localhost:8082/v1/tasks/search")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.String())

	return resp.Body()
}

// GetAllCreatedTasks
// @Description: 獲得所有狀態為 Created 任務詳細資訊
func GetAllCreatedTasks() []byte {
	client := resty.New()

	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"state": "CREATED",
		}).
		Post("http://localhost:8082/v1/tasks/search")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.String())

	return resp.Body()
}

// GetSpecificTask
// @Description: 獲得特定 ID 任務詳細資訊
func GetSpecificTask(taskId string) []byte {
	client := resty.New()

	resp, err := client.R().SetHeader("Content-Type", "application/json").
		Get("http://localhost:8082/v1/tasks/" + taskId)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.String())

	return resp.Body()
}

// AssignTask
// @Description: 指派任務給使用者
func AssignTask(taskId string, userName string) []byte {
	client := resty.New()

	resp, err := client.R().SetHeader("Content-Type", "application/json").SetHeader("Authorization", token).
		SetBody(map[string]interface{}{
			"assignee":                userName,
			"allowOverrideAssignment": true,
		}).
		Patch("http://localhost:8082/v1/tasks/" + taskId + "/assign")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.String())

	return resp.Body()
}

// Unassigntask
// @Description: 取消指派任務
func Unassigntask(taskId string) []byte {
	client := resty.New()
	name, value := GetCookie()

	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetCookie(&http.Cookie{Name: name, Value: value}).
		Patch("http://localhost:8082/v1/tasks/" + taskId + "/unassign")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.String())

	return resp.Body()
}

// CompleteTask
// @Description: 完成任務
func CompleteTask(taskId string, variables []map[string]string) []byte {
	client := resty.New()
	name, value := GetCookie()

	reqBody := map[string]interface{}{
		"variables": variables,
	}

	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetCookie(&http.Cookie{Name: name, Value: value}).
		SetBody(reqBody).
		Patch("http://localhost:8082/v1/tasks/" + taskId + "/complete")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.String())

	return resp.Body()
}

// GetTaskVariables
// @Description: 獲得任務的變數清單
func GetTaskVariables(taskId string) []byte {
	client := resty.New()
	name, value := GetCookie()

	resp, err := client.R().SetHeader("accept", "application/json").
		SetCookie(&http.Cookie{Name: name, Value: value}).
		Post("http://localhost:8082/v1/tasks/" + taskId + "/variables/search")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.String())

	return resp.Body()
}
