package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

func main() {

	AssignTask("2251799813693381", "demo")
	variables := []map[string]string{
		{
			"name":  "item",
			"value": "\"Task Assignment\"",
		},
	}
	CompleteTask("2251799813693381", variables)

	//GetTaskVariables("2251799813693348")

	//GetAllCreatedTasks()
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
	name, value := GetCookie()

	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetCookie(&http.Cookie{Name: name, Value: value}).
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
	name, value := GetCookie()

	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetCookie(&http.Cookie{Name: name, Value: value}).
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
	name, value := GetCookie()

	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetCookie(&http.Cookie{Name: name, Value: value}).
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
	name, value := GetCookie()

	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetCookie(&http.Cookie{Name: name, Value: value}).
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
